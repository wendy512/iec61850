package scl

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/cast"
	"gopkg.in/validator.v2"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	filePath string
}

func NewParser(filePath string) *Parser {
	return &Parser{
		filePath: filePath,
	}
}

func (p *Parser) Parse() (*SCL, error) {
	fileBytes, err := os.ReadFile(p.filePath)

	if err != nil {
		return nil, err
	}

	var scl SCL
	if err = xml.Unmarshal(fileBytes, &scl); err != nil {
		return nil, err
	}
	scl.FileFullPath = p.filePath

	if err = validator.Validate(scl); err != nil {
		return nil, err
	}

	typeDeclarations := make([]SclType, 0)

	if scl.DataTypeTemplates.DataAttributeTypes != nil {

		for _, typeDeclare := range scl.DataTypeTemplates.DataAttributeTypes {

			// set default values
			if typeDeclare.SubDataAttributes == nil {
				typeDeclare.SubDataAttributes = make([]*DataAttributeDefinition, 0)
			}

			if err = p.parseDataAttributeDefinition(typeDeclare.SubDataAttributes); err != nil {
				return nil, err
			}

			typeDeclarations = append(typeDeclarations, typeDeclare)
		}
	}

	if scl.DataTypeTemplates.DataObjectTypes != nil {

		for _, typeDeclare := range scl.DataTypeTemplates.DataObjectTypes {

			// set default values
			if typeDeclare.DataAttributes == nil {
				typeDeclare.DataAttributes = make([]*DataAttributeDefinition, 0)
			}

			if err = p.parseDataAttributeDefinition(typeDeclare.DataAttributes); err != nil {
				return nil, err
			}

			typeDeclarations = append(typeDeclarations, typeDeclare)
		}
	}

	if scl.DataTypeTemplates.EnumTypes != nil {
		for _, typeDeclare := range scl.DataTypeTemplates.EnumTypes {
			typeDeclarations = append(typeDeclarations, typeDeclare)
		}
	}

	if scl.DataTypeTemplates.LogicalNodeTypes != nil {
		for _, typeDeclare := range scl.DataTypeTemplates.LogicalNodeTypes {
			typeDeclarations = append(typeDeclarations, typeDeclare)
		}
	}

	scl.DataTypeTemplates.TypeDeclarations = typeDeclarations

	for _, ied := range scl.IEDs {

		for _, ap := range ied.AccessPoints {

			for _, lDevice := range ap.Server.LogicalDevices {

				logicNodes := make([]*LogicalNode, 0)

				if lDevice.LN0 != nil {

					lnType := lDevice.LN0.LnType
					_sclType := scl.DataTypeTemplates.lookupType(lnType)
					if _sclType == nil {
						return nil, fmt.Errorf("%s missing type declaration %s", lDevice.LN0.GetName(), lnType)
					}

					lDevice.LN0.SclType = _sclType
					logicNodes = append(logicNodes, lDevice.LN0)
				}

				if lDevice.LNodes != nil {

					for _, ln := range lDevice.LNodes {

						lnType := ln.LnType
						_sclType := scl.DataTypeTemplates.lookupType(lnType)
						if _sclType == nil {
							return nil, fmt.Errorf("%s missing type declaration %s", ln.GetName(), lnType)
						}

						ln.SclType = _sclType
						logicNodes = append(logicNodes, ln)
					}
				}

				for _, logicNode := range logicNodes {

					if lNodeType, ok := logicNode.SclType.(*LogicalNodeType); ok {
						// mark type as used
						lNodeType.SetUsed(true)

						logicNode.DataObjects = make([]*DataObject, 0)
						for _, doDefinition := range lNodeType.DataObjectDefinitions {

							da, err := NewDataObject(doDefinition, scl.DataTypeTemplates, logicNode)
							if err != nil {
								return nil, err
							}
							logicNode.DataObjects = append(logicNode.DataObjects, da)
						}

						// set default values for control blocks
						if logicNode.DataSets == nil {
							logicNode.DataSets = make([]*DataSet, 0)
						}
						if logicNode.ReportControlBlocks == nil {
							logicNode.ReportControlBlocks = make([]*ReportControl, 0)
						}
						if logicNode.GSEControlBlocks == nil {
							logicNode.GSEControlBlocks = make([]*GSEControl, 0)
						}
						if logicNode.SMVControlBlocks == nil {
							logicNode.SMVControlBlocks = make([]*SampledValueControl, 0)
						}
						if logicNode.LogControlBlocks == nil {
							logicNode.LogControlBlocks = make([]*LogControl, 0)
						}
						if logicNode.Logs == nil {
							logicNode.Logs = make([]*Log, 0)
						}
						if logicNode.SettingGroupControlBlocks == nil {
							logicNode.SettingGroupControlBlocks = make([]*SettingControl, 0)
						}
						if logicNode.DOINodes == nil {
							logicNode.DOINodes = make([]*DOINode, 0)
						}

						for _, rcb := range logicNode.ReportControlBlocks {

							if rcb.TriggerOptions == nil {

								// use default values if no node present
								rcb.TriggerOptions = &TriggerOptions{
									Gi: true,
								}
							}

							rcb.Indexed = true
							if rcb.IndexedStr != "" {
								rcb.Indexed = cast.ToBool(rcb.IndexedStr)
							}

							if rcb.RptEnabled != nil {
								// set default value
								rcb.RptEnabled.Max = 1

								if rcb.RptEnabled.MaxStr != "" {
									if rcb.RptEnabled.Max, err = cast.ToIntE(rcb.RptEnabled.MaxStr); err != nil {
										return nil, err
									}
								}

								if !rcb.Indexed {
									if rcb.RptEnabled.Max != 1 {
										return nil, fmt.Errorf("rptEnabled.max != 1 is not allowed when indexed=\"false\"")
									}
								}
							}
						}

						for _, scb := range logicNode.GSEControlBlocks {
							if scb.Type != "" && scb.Type != "GOOSE" {
								return nil, fmt.Errorf("GSEControl of type %s not supported!", scb.Type)
							}
						}

						for _, scb := range logicNode.SMVControlBlocks {
							switch scb.SmpModStr {
							case "SmpPerPeriod":
								scb.SmpMod = SmpPerPeriod
							case "SmpPerSec":
								scb.SmpMod = SmpPerSecond
							case "SecPerSmp":
								scb.SmpMod = SecPerSmp
							default:
								return nil, fmt.Errorf("invalid smpMod value %s", scb.SmpModStr)
							}
						}

						for _, lcb := range logicNode.LogControlBlocks {
							if lcb.TriggerOptions == nil {

								// use default values if no node present
								lcb.TriggerOptions = &TriggerOptions{
									Gi: true,
								}
							}
						}

						for _, log := range logicNode.Logs {
							if log.Name == "" {
								// set default value
								log.Name = "GeneralLog"
							}
						}

						if !(logicNode.LnClass == "LLN0") && len(logicNode.SettingGroupControlBlocks) > 0 {
							return nil, fmt.Errorf("LN other than LN0 is not allowed to contain SettingControl")
						}

						for _, doiNode := range logicNode.DOINodes {
							dataModelNode := logicNode.GetChildByName(doiNode.Name)

							if dataModelNode == nil {
								return nil, fmt.Errorf("missing data object with name \"%s\"", doiNode.Name)
							}

							if err = p.parseDataAttributeNodes(doiNode, dataModelNode); err != nil {
								return nil, err
							}
						}

					} else {
						return nil, fmt.Errorf("wrong type %s for logical node", logicNode.LnType)
					}
				}

				lDevice.LogicalNodes = logicNodes
			}
		}
	}

	if scl.Communication != nil {

		// set default values for communication
		if scl.Communication.SubNetworks == nil {
			scl.Communication.SubNetworks = make([]*SubNetwork, 0)
		}

		for _, subNetwork := range scl.Communication.SubNetworks {
			if subNetwork.ConnectedAP == nil {
				subNetwork.ConnectedAP = make([]*ConnectedAP, 0)
			}

			for _, connectedAP := range subNetwork.ConnectedAP {
				if connectedAP.GESNodes == nil {
					connectedAP.GESNodes = make([]*GSE, 0)
				}

				if connectedAP.SMVNodes == nil {
					connectedAP.SMVNodes = make([]*SMV, 0)
				}

				for _, node := range connectedAP.GESNodes {
					node.MaxTime = -1
					node.MinTime = -1

					if node.MaxTimeVal != nil {
						if node.MaxTime, err = cast.ToIntE(node.MaxTimeVal.Value); err != nil {
							return nil, err
						}
					}

					if node.MinTimeVal != nil {
						if node.MinTime, err = cast.ToIntE(node.MinTimeVal.Value); err != nil {
							return nil, err
						}
					}

					if err = p.parsePhyComAddress(node.Address); err != nil {
						return nil, err
					}
				}

				for _, node := range connectedAP.SMVNodes {
					if err = p.parsePhyComAddress(node.Address); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return &scl, nil
}

func (p *Parser) parsePhyComAddress(phyComAddress *PhyComAddress) error {
	var err error

	if phyComAddress != nil && phyComAddress.AddressParameters != nil {
		phyComAddress.VlanId = -1
		phyComAddress.VlanPriority = -1
		phyComAddress.AppId = -1

		for _, p2 := range phyComAddress.AddressParameters {

			switch p2.Type {
			case "VLAN-ID":
				phyComAddress.VlanId, err = strconv.ParseInt(p2.Value, 16, 32)
				if err != nil || phyComAddress.VlanId > 0xfff {
					return fmt.Errorf("VLAN-ID value out of range: %s", p2.Value)
				}
			case "VLAN-PRIORITY":
				phyComAddress.VlanPriority, err = strconv.Atoi(p2.Value)
				if err != nil {
					return fmt.Errorf("invalid VLAN-PRIORITY: %s", p2.Value)
				}
			case "APPID":
				phyComAddress.AppId, err = strconv.ParseInt(p2.Value, 16, 32)
				if err != nil || phyComAddress.AppId > 0xffff {
					return fmt.Errorf("APPID value out of range: %s", p2.Value)
				}
			case "MAC-Address":
				addressElements := strings.Split(p2.Value, "-")
				if len(addressElements) != 6 {
					return fmt.Errorf("malformed address: %s", p2.Value)
				}

				phyComAddress.MacAddress = make([]int, 6)
				for i, element := range addressElements {
					parsed, err := strconv.ParseUint(element, 16, 8)
					if err != nil {
						return fmt.Errorf("invalid MAC address element: %s", element)
					}
					phyComAddress.MacAddress[i] = int(parsed)
				}
			}
		}

		if phyComAddress.VlanId == -1 {
			phyComAddress.VlanId = 0
		}

		if phyComAddress.VlanPriority == -1 {
			phyComAddress.VlanPriority = 4
		}

		if phyComAddress.AppId == -1 {
			phyComAddress.AppId = 0
		}

		if phyComAddress.MacAddress == nil {
			phyComAddress.MacAddress = []int{0x01, 0x0C, 0xCD, 0x01, 0x00, 0x00}
		}
	}
	return nil
}

func (p *Parser) parseDataAttributeDefinition(daDefinitions []*DataAttributeDefinition) error {
	for _, subDataAttribute := range daDefinitions {

		if attributeType, err := createAttributeTypeFromScl(subDataAttribute.BType); err != nil {
			return err
		} else {

			switch subDataAttribute.BType {
			case "Tcmd":
				subDataAttribute.Type = "Tcmd"
			case "Dbpos":
				subDataAttribute.Type = "Dbpos"
			case "Check":
				subDataAttribute.Type = "Check"
			}

			subDataAttribute.AttributeType = attributeType
			subDataAttribute.TriggerOptions = &TriggerOptions{
				Dchg: subDataAttribute.DchgTrigger,
				Qchg: subDataAttribute.QchgTrigger,
				Dupd: subDataAttribute.DupdTrigger,
			}

			if subDataAttribute.BType != "" {

				if subDataAttribute.Val != nil {

					if attributeType == Enumerated {

						subDataAttribute.Value = &DataModelValue{
							UnknownEnumValue: subDataAttribute.Val.Value,
							EnumType:         subDataAttribute.Type,
						}
					} else {

						dataModelValue, err := NewDataModelValue(attributeType, nil, subDataAttribute.Val.Value)
						if err != nil {
							return err
						}
						subDataAttribute.Value = dataModelValue
					}
				}

			}
		}

	}
	return nil
}

func (p *Parser) parseDataAttributeNodes(doiNode *DOINode, do DataModelNode) error {
	if doiNode.DAINodes == nil {
		doiNode.DAINodes = make([]*DOINode, 0)
	}

	for _, daiNode := range doiNode.DAINodes {
		dataModelNode := do.GetChildByName(daiNode.Name)

		if dataModelNode == nil {
			return fmt.Errorf("missing data attribute with name \"%s\"", daiNode.Name)
		}

		dataAttribute := dataModelNode.(*DataAttribute)
		if val := daiNode.Val; val != nil {
			dataModelValue, err := NewDataModelValue(dataAttribute.AttributeType, dataAttribute.SclType, val.Value)

			if err != nil {
				return err
			}
			dataAttribute.Value = dataModelValue

			if daiNode.SAddr != "" {
				dataAttribute.ShortAddress = daiNode.SAddr
			}
		}
	}

	return nil
}

func (p *Parser) parseSubDataInstances(doiNode *DOINode, do DataModelNode) error {
	if doiNode.SDINodes == nil {
		doiNode.SDINodes = make([]*DOINode, 0)
	}

	for _, sdiNode := range doiNode.SDINodes {

		dataModelNode := do.GetChildByName(sdiNode.Name)
		if dataModelNode == nil {
			return fmt.Errorf("subelement with name %s not found", sdiNode.Name)
		}

		if err := p.parseDataAttributeNodes(sdiNode, dataModelNode); err != nil {
			return err
		}

		if err := p.parseSubDataInstances(sdiNode, dataModelNode); err != nil {
			return err
		}
	}

	return nil
}
