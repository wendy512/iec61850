package scl_xml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// IEC 61850 SCL (ICD) Data Structures

// DAValue interface for all possible data types
type DAValue interface{}

// Simple Data Types in Go for IEC 61850

type BOOLEAN bool
type INT32 int32
type INT16 int16
type INT8 int8
type FLOAT32 float32
type VisString255 string
type Unicode255 string

type DataSetDetail struct {
	DataSet
	IEDName           string
	DOTypes           map[string]DOType
	DataTypeTemplates DataTypeTemplates
}

func (ds *DataSetDetail) GetDOType(prefix, lnClass, name string) DOType {
	var doTypeId string

	var builder strings.Builder
	builder.WriteString(prefix)
	builder.WriteString(lnClass)
	builder.WriteString("/")
	builder.WriteString(name)

	defer builder.Reset()

	if ds.DOTypes != nil {
		if typ, ok := ds.DOTypes[builder.String()]; ok {
			return typ
		}
	}

	for _, lNodeType := range ds.DataTypeTemplates.LNodeType {
		if lNodeType.ID == fmt.Sprintf("%s%s", prefix, lnClass) || lNodeType.LNClass == lnClass {
			for _, do := range lNodeType.DO {
				if do.Name == name {
					doTypeId = do.Type
					break
				}
			}
			break
		}
	}

	for _, doType := range ds.DataTypeTemplates.DOType {
		if doType.ID == doTypeId {
			if ds.DOTypes == nil {
				ds.DOTypes = make(map[string]DOType)
			}

			ds.DOTypes[builder.String()] = doType
			return doType
		}
	}

	return DOType{}
}

type SCL struct {
	IED               []IED             `xml:"IED"`
	DataTypeTemplates DataTypeTemplates `xml:"DataTypeTemplates"`
}

func (scl *SCL) GetDataSet(ref string) (*DataSetDetail, error) {
	args := strings.Split(ref, "/")
	if len(args) != 2 {
		return nil, fmt.Errorf("error parse dataset ref: %s", ref)
	}

	for _, ied := range scl.IED {
		for _, accessPoint := range ied.AccessPoint {
			for _, lDevice := range accessPoint.LDevice {
				if fmt.Sprintf("%s%s", ied.Name, lDevice.Inst) == args[0] {
					for _, dSet := range lDevice.LN0.DataSets {
						if fmt.Sprintf("%s.%s", lDevice.LN0.LnClass, dSet.Name) == args[1] {
							return &DataSetDetail{
								IEDName:           ied.Name,
								DataSet:           dSet,
								DataTypeTemplates: scl.DataTypeTemplates,
							}, nil
						}
					}
					break
				}
			}
		}
	}

	return nil, fmt.Errorf("can not found dataset ref: %s", ref)
}

type IED struct {
	Name          string        `xml:"name,attr"`
	Type          string        `xml:"type,attr"`
	Desc          string        `xml:"desc,attr"`
	ConfigVersion string        `xml:"configVersion,attr"`
	AccessPoint   []AccessPoint `xml:"AccessPoint"`
}

type AccessPoint struct {
	Name    string    `xml:"name,attr"`
	LDevice []LDevice `xml:"Server>LDevice"`
}

type LDevice struct {
	Inst string `xml:"inst,attr"`
	LN   []LN   `xml:"LN"`
	LN0  LN0    `xml:"LN0"`
}

type LN0 struct {
	Inst     string    `xml:"inst,attr"`
	LnType   string    `xml:"lnType,attr"`
	LnClass  string    `xml:"lnClass,attr"`
	DataSets []DataSet `xml:"DataSet"`
}

type LN struct {
	Inst    string `xml:"inst,attr"`
	Prefix  string `xml:"prefix,attr"`
	LnType  string `xml:"lnType,attr"`
	LnClass string `xml:"lnClass,attr"`
	DOI     []DOI  `xml:"DOI"`
}

type DOI struct {
	Desc string `xml:"desc,attr"`
	Name string `xml:"name,attr"`
	DAI  []DAI  `xml:"DAI"`
	SDI  []SDI  `xml:"SDI"`
}

type DAI struct {
	Name string `xml:"name,attr"`
	Val  Val    `xml:"Val"`
	SDI  []SDI  `xml:"SDI"` // 新增
}

type SDI struct {
	Name string `xml:"name,attr"`
	DAI  []DAI  `xml:"DAI"`
	SDI  []SDI  `xml:"SDI"` // 递归地包含SDI
}

type Val struct {
	Value string `xml:",chardata"`
}

type DataSet struct {
	Name string      `xml:"name,attr"`
	Desc string      `xml:"desc,attr"`
	FCDA []FCDAEntry `xml:"FCDA"`
}

type FCDAEntry struct {
	LDInst  string `xml:"ldInst,attr,omitempty"`
	Prefix  string `xml:"prefix,attr,omitempty"`
	LNClass string `xml:"lnClass,attr"`
	LNInst  string `xml:"lnInst,attr,omitempty"`
	DOName  string `xml:"doName,attr"`
	DAName  string `xml:"daName,attr,omitempty"`
	FC      string `xml:"fc,attr"`
}

type DataTypeTemplates struct {
	LNodeType []LNodeType `xml:"LNodeType"`
	DOType    []DOType    `xml:"DOType"`
	DAType    []DAType    `xml:"DAType"`
	EnumType  []EnumType  `xml:"EnumType"`
}

type LNodeType struct {
	ID      string `xml:"id,attr"`
	LNClass string `xml:"lnClass,attr"`
	Desc    string `xml:"desc,attr,omitempty"`
	DO      []DO   `xml:"DO"`
}

type DO struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr,omitempty"`
	DA   []DA   `xml:"DA"`
}

type DOType struct {
	ID   string `xml:"id,attr"`
	DA   []DA   `xml:"DA"`
	Desc string `xml:"desc,attr,omitempty"`
	SDO  []SDO  `xml:"SDO"`
}

type DA struct {
	Name string  `xml:"name,attr"`
	Type string  `xml:"bType,attr"`
	FC   string  `xml:"fc,attr"`
	Val  DAValue `xml:"Val"`
	DA   []DA    `xml:"DA"`
}

type DAType struct {
	ID  string `xml:"id,attr"`
	BDA []BDA  `xml:"BDA"`
	DA  []DA   `xml:"DA"`
}

type BDA struct {
	Name string  `xml:"name,attr"`
	Type string  `xml:"type,attr"`
	Val  DAValue `xml:"Val"`
}

type EnumType struct {
	ID      string    `xml:"id,attr"`
	EnumVal []EnumVal `xml:"EnumVal"`
}

type EnumVal struct {
	Ord  int    `xml:"ord,attr"`
	Name string `xml:",chardata"`
}

type SDO struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	DA   []DA   `xml:"DA"`
}

func (scl *SCL) Print() {
	for _, ied := range scl.IED {
		ied.Print(0)
	}
	scl.DataTypeTemplates.Print(0)
}

func (ied *IED) Print(depth int) {
	fmt.Printf("%sIED Name: %s, Type: %s, Desc: %s\n", getIndentation(depth), ied.Name, ied.Type, ied.Desc)
	for _, ap := range ied.AccessPoint {
		fmt.Printf("%sAccessPoint: %s\n", getIndentation(depth+1), ap.Name)
		for _, ld := range ap.LDevice {
			fmt.Printf("%sLDevice: %s\n", getIndentation(depth+2), ld.Inst)
			for _, ln := range ld.LN {
				ln.Print(depth + 3)
			}
		}
	}
}

func (ln *LN) Print(depth int) {
	fmt.Printf("%sLN Inst: %s, Prefix: %s, LnType: %s, LnClass: %s\n", getIndentation(depth), ln.Inst, ln.Prefix, ln.LnType, ln.LnClass)
	for _, doi := range ln.DOI {
		doi.Print(depth + 1)
	}
}

func (doi *DOI) Print(depth int) {
	fmt.Printf("%sDOI Name: %s, Desc: %s\n", getIndentation(depth), doi.Name, doi.Desc)
	for _, dai := range doi.DAI {
		dai.Print(depth + 1)
	}
	// Here, assuming you also want to print SDIs if they are included in your model
	for _, sdi := range doi.SDI {
		sdi.Print(depth + 1)
	}
}

func (dai *DAI) Print(depth int) {
	fmt.Printf("%sDAI Name: %s, Value: %s\n", getIndentation(depth), dai.Name, dai.Val.Value)
	for _, sdi := range dai.SDI {
		sdi.Print(depth + 1)
	}
}

func (sdi *SDI) Print(depth int) {
	// Print SDI and its related DAIs
	// Note: If SDIs can have nested SDIs, this function will need recursion
	fmt.Printf("%sSDI: %s\n", getIndentation(depth), sdi.Name)
	for _, dai := range sdi.DAI {
		fmt.Printf("%sDAI: %s\n", getIndentation(depth+1), dai.Name)
		dai.Val.Print(depth + 2)
	}
}

func (val *Val) Print(depth int) {
	fmt.Printf("%sValue: %s\n", getIndentation(depth), val.Value)
}

func (dt DataTypeTemplates) Print(depth int) {
	fmt.Printf("%sDataTypeTemplates:\n", getIndentation(depth))
	for _, lnt := range dt.LNodeType {
		fmt.Printf("%sLNodeType: %s\n", getIndentation(depth+1), lnt.ID)
		for _, dt := range lnt.DO {
			dt.Print(depth + 2)
		}
	}
	for _, dot := range dt.DOType {
		fmt.Printf("%sDOType: %s\n", getIndentation(depth+1), dot.ID)
		for _, da := range dot.DA {
			da.Print(depth + 2)
		}
	}
	for _, dat := range dt.DAType {
		fmt.Printf("%sDAType: %s\n", getIndentation(depth+1), dat.ID)
		for _, bda := range dat.BDA {
			bda.Print(depth + 2)
		}
		for _, da := range dat.DA {
			da.Print(depth + 2)
		}
	}
	for _, et := range dt.EnumType {
		fmt.Printf("%sEnumType: %s\n", getIndentation(depth+1), et.ID)
		for _, ev := range et.EnumVal {
			fmt.Printf("%sEnumVal: Ord: %d, Name: %s\n", getIndentation(depth+2), ev.Ord, ev.Name)
		}
	}
}

func (do DO) Print(depth int) {
	if do.Name != "" && do.Type != "" {
		fmt.Printf("%sDO: %s, Type: %s\n", getIndentation(depth), do.Name, do.Type)
	}
}

func (da DA) Print(depth int) {
	if da.Name != "" && da.Type != "" {
		fmt.Printf("%sDA: %s, Type: %s\n", getIndentation(depth), da.Name, da.Type)
		for _, subDA := range da.DA {
			subDA.Print(depth + 1)
		}
	}
}

func (bda BDA) Print(depth int) {
	if bda.Name != "" && bda.Type != "" {
		fmt.Printf("%sBDA: %s, Type: %s\n", getIndentation(depth), bda.Name, bda.Type)
	}
}

func getIndentation(depth int) string {
	return strings.Repeat("  ", depth)
}

func GetSCL(path string) (SCL, error) {
	var scl SCL
	// 打开并读取ICD文件
	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return scl, fmt.Errorf("open file failed: %v", err)
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return scl, err
	}

	err = xml.Unmarshal(byteValue, &scl)
	if err != nil {
		return scl, fmt.Errorf("unmarshall failed: %v", err)
	}

	return scl, nil
}
