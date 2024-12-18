package scl

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cast"
	"os"
	"strconv"
	"strings"
	"time"
)

type PrintStream struct {
	builder *strings.Builder
}

func (p *PrintStream) println(format string, a ...any) {
	_, _ = fmt.Fprintf(p.builder, format+"\n", a...)
}

func (p *PrintStream) printlnNone() {
	_, _ = fmt.Fprintln(p.builder)
}

func (p *PrintStream) print(format string, a ...any) {
	_, _ = fmt.Fprintf(p.builder, format, a...)
}

func (p *PrintStream) writeFile(fileFullPath string) error {
	return os.WriteFile(fileFullPath, []byte(p.builder.String()), 0644)
}

type DataObject struct {
	Name              string
	Count             int
	Trans             bool
	SclType           SclType
	DataTypeTemplates *DataTypeTemplates
	DataAttributes    []*DataAttribute
	SubDataObjects    []*DataObject
}

func NewDataObject(do *DataObjectDefinition, dataTypeTemplates *DataTypeTemplates, parent DataModelNode) (*DataObject, error) {
	_sclType := dataTypeTemplates.lookupType(do.Type)
	if _sclType == nil {
		return nil, fmt.Errorf("%s missing type declaration %s", do.Name, do.Type)
	}

	_sclType.SetUsed(true)
	dataObject := &DataObject{
		Name:              do.Name,
		Count:             do.Count,
		Trans:             do.Transient,
		DataTypeTemplates: dataTypeTemplates,
		SclType:           _sclType,
		DataAttributes:    make([]*DataAttribute, 0),
		SubDataObjects:    make([]*DataObject, 0),
	}

	if err := dataObject.createDataAttributes(); err != nil {
		return nil, err
	}

	if err := dataObject.createSubDataObjects(); err != nil {
		return nil, err
	}

	return dataObject, nil
}

func (d *DataObject) GetName() string {
	return d.Name
}

func (d *DataObject) GetSclType() SclType {
	return d.SclType
}

func (d *DataObject) GetChildByName(childName string) DataModelNode {

	if d.DataAttributes != nil {
		for _, dataAttribute := range d.DataAttributes {
			if dataAttribute.GetName() == childName {
				return dataAttribute
			}
		}
	}

	if d.SubDataObjects != nil {
		for _, subDataObject := range d.SubDataObjects {
			if subDataObject.GetName() == childName {
				return subDataObject
			}
		}
	}

	return nil
}

func (d *DataObject) createDataAttributes() error {
	var daDefinitions []*DataAttributeDefinition

	if doType, ok := d.SclType.(*DataObjectType); ok {
		daDefinitions = doType.DataAttributes
	}

	if daType, ok := d.SclType.(*DataAttributeType); ok {
		daDefinitions = daType.SubDataAttributes
	}

	for _, daDefinition := range daDefinitions {
		if daDefinition.Fc == "SE" {

			attribute, err := NewDataAttribute(daDefinition, d.DataTypeTemplates, "SG", d)
			if err != nil {
				return err
			}
			d.DataAttributes = append(d.DataAttributes, attribute)
		}

		attribute, err := NewDataAttribute(daDefinition, d.DataTypeTemplates, "", d)
		if err != nil {
			return err
		}

		isMulti := false
		for i, da := range d.DataAttributes {
			if da.FC != "SG" && da.GetName() == attribute.GetName() {

				fmt.Printf("Warning: DataObject %s has multi attribute %s, fc %s\n", d.GetName(), attribute.GetName(), da.FC)
				if da.FC == "SP" && attribute.FC == "SE" {
					isMulti = true
					d.DataAttributes[i] = attribute
					break
				}
			}
		}

		if !isMulti {
			d.DataAttributes = append(d.DataAttributes, attribute)
		}
	}

	return nil
}

func (d *DataObject) createSubDataObjects() error {
	for _, sdoDefinition := range d.SclType.(*DataObjectType).SubDataObjects {
		dataObject, err := NewDataObject(sdoDefinition, d.DataTypeTemplates, d)

		if err != nil {
			return err
		}

		d.SubDataObjects = append(d.SubDataObjects, dataObject)
	}

	return nil
}

type DataAttribute struct {
	Name              string
	FC                string
	AttributeType     AttributeType
	IsBasicAttribute  bool
	Count             int
	Value             *DataModelValue
	ShortAddress      string
	SubDataAttributes []*DataAttribute
	SclType           SclType
	TriggerOptions    *TriggerOptions
	Definition        *DataAttributeDefinition
}

func NewDataAttribute(daDefinition *DataAttributeDefinition, dataTypeTemplates *DataTypeTemplates, fc string, parent DataModelNode) (*DataAttribute, error) {
	da := &DataAttribute{
		Name:          daDefinition.Name,
		FC:            daDefinition.Fc,
		AttributeType: daDefinition.AttributeType,
		Count:         daDefinition.Count,
		Definition:    daDefinition,
	}

	if da.FC == "" {
		da.FC = fc
	}

	if fc != "" {
		da.FC = fc
	}

	if parent != nil {

		if node, ok := parent.(*DataAttribute); ok {
			da.TriggerOptions = node.TriggerOptions
		} else {
			da.TriggerOptions = daDefinition.TriggerOptions
		}

	} else {
		da.TriggerOptions = daDefinition.TriggerOptions
	}

	if da.AttributeType == Constructed {
		da.IsBasicAttribute = false
		if err := da.createConstructedAttribute(dataTypeTemplates); err != nil {
			return nil, err
		}
	} else if da.AttributeType == Enumerated {

		if err := da.createEnumeratedAttribute(dataTypeTemplates); err != nil {
			return nil, err
		}
	}

	return da, nil
}

func (d *DataAttribute) GetName() string {
	return d.Name
}

func (d *DataAttribute) GetSclType() SclType {
	return d.SclType
}

func (d *DataAttribute) GetChildByName(childName string) DataModelNode {
	if d.SubDataAttributes != nil {
		for _, subDataAttribute := range d.SubDataAttributes {
			if subDataAttribute.GetName() == childName {
				return subDataAttribute
			}
		}
	}

	return nil
}

func (d *DataAttribute) createEnumeratedAttribute(dataTypeTemplates *DataTypeTemplates) error {
	_sclType := dataTypeTemplates.lookupType(d.Definition.Type)
	if _sclType == nil {
		return fmt.Errorf("missing type definition for enumerated data attribute: %s", d.Definition.Type)
	}
	d.SclType = _sclType

	if _, ok := _sclType.(*EnumerationType); !ok {
		return fmt.Errorf("wrong type definition for enumerated data attribute")
	}

	_sclType.SetUsed(true)
	return nil
}

func (d *DataAttribute) createConstructedAttribute(dataTypeTemplates *DataTypeTemplates) error {
	_sclType := dataTypeTemplates.lookupType(d.Definition.Type)
	if _sclType == nil {
		return fmt.Errorf("missing type definition for constructed data attribute: %s", d.Definition.Type)
	}
	d.SclType = _sclType

	if dataAttributeType, ok := _sclType.(*DataAttributeType); !ok {
		return fmt.Errorf("wrong type definition for constructed data attribute")
	} else {

		_sclType.SetUsed(true)
		d.SubDataAttributes = make([]*DataAttribute, 0)

		for _, daDef := range dataAttributeType.SubDataAttributes {
			dataAttribute, err := NewDataAttribute(daDef, dataTypeTemplates, d.FC, d)
			if err != nil {
				return err
			}
			d.SubDataAttributes = append(d.SubDataAttributes, dataAttribute)
		}
	}

	return nil
}

type DataModelValue struct {
	Value            interface{}
	UnknownEnumValue string
	EnumType         string
}

func NewDataModelValue(attributeType AttributeType, sclType SclType, value string) (*DataModelValue, error) {
	dmv := &DataModelValue{}
	var err error

	switch attributeType {
	case Enumerated:
		enumType := sclType.(*EnumerationType)
		var ord int

		if ord, err = enumType.getOrdByEnumString(value); err == nil {
			dmv.Value = ord
		} else {
			if ord, err = strconv.Atoi(value); err == nil {
				if enumType.isValidOrdValue(ord) {
					dmv.Value = ord
				} else {
					return nil, fmt.Errorf("%s is not a valid value of type %s", value, sclType.GetId())
				}
			} else {
				return nil, fmt.Errorf("%s is not a valid value of type %s", value, sclType.GetId())
			}
		}
	case Int8, Int16, Int32, Int8U, Int16U, Int24U, Int32U, Int64:
		if trimmedValue := strings.TrimSpace(value); trimmedValue == "" {
			dmv.Value = int64(0)
		} else {
			dmv.Value, err = cast.ToInt64E(trimmedValue)
		}
	case Boolean:
		dmv.Value = strings.ToLower(strings.TrimSpace(value)) == "true"
	case Float32:
		trimmedValue := strings.ReplaceAll(strings.TrimSpace(value), ",", ".")
		if trimmedValue == "" {
			dmv.Value = float32(0)
		} else {
			dmv.Value, err = cast.ToFloat32E(trimmedValue)
		}
	case Float64:
		trimmedValue := strings.ReplaceAll(strings.TrimSpace(value), ",", ".")
		if trimmedValue == "" {
			dmv.Value = float64(0)
		} else {
			dmv.Value, err = cast.ToFloat64E(trimmedValue)
		}
	case UnicodeString255, VisibleString32, VisibleString64, VisibleString65, VisibleString129, VisibleString255, Currency:
		dmv.Value = value
	case OctetString64:
		dmv.Value, err = base64.StdEncoding.DecodeString(value)
	case Check:
		fmt.Println("Warning: Initialization of CHECK is unsupported!")
	case CodedEnum:
		switch value {
		case "intermediate-state", "stop":
			dmv.Value = 0
		case "off", "lower":
			dmv.Value = 1
		case "on", "higher":
			dmv.Value = 2
		case "bad-state", "reserved":
			dmv.Value = 4
		default:
			fmt.Printf("Warning: CODEDENUM is initialized with unsupported value %s\n", value)
		}
	case Quality:
		fmt.Println("Warning: Initialization of QUALITY is unsupported!")
	case Timestamp, EntryTime:
		modValueString := strings.ReplaceAll(value, ",", ".")
		date, err := time.Parse("2006-01-02T15:04:05.000", modValueString)
		if err != nil {
			fmt.Printf("Warning: Val element does not contain a valid time stamp: %s\n", err)
		} else {
			dmv.Value = date.UnixMilli()
		}
	default:
		return nil, fmt.Errorf("unsupported type %d value: %s", attributeType, value)
	}

	if err != nil {
		return nil, err
	}
	return dmv, err
}

func (that *DataModelValue) updateEnumOrdValue(templates *DataTypeTemplates) {
	if that.EnumType != "" {
		if _sclType := templates.lookupType(that.EnumType); _sclType != nil {

			enumType := _sclType.(*EnumerationType)
			if ord, err := enumType.getOrdByEnumString(that.UnknownEnumValue); err == nil {
				that.Value = ord
			} else {
				fmt.Printf("failed: %s\n", err)
			}
		}
	}
}
