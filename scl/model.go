package scl

import (
	"encoding/xml"
	"fmt"
)

type SCL struct {
	XMLName           xml.Name           `xml:"SCL"`
	Header            *Header            `xml:"Header"`
	Communication     *Communication     `xml:"Communication"`
	IEDs              []*IED             `xml:"IED" validate:"nonzero"`
	DataTypeTemplates *DataTypeTemplates `xml:"DataTypeTemplates" validate:"nonzero"`

	// custom
	FileFullPath string
}

type Header struct {
	ID            string `xml:"id,attr"`
	ToolID        string `xml:"toolID,attr"`
	NameStructure string `xml:"nameStructure,attr"`
}

type Communication struct {
	SubNetworks []*SubNetwork `xml:"SubNetwork" validate:"nonzero"`
}

type SubNetwork struct {
	Name        string         `xml:"name,attr"`
	Type        string         `xml:"type,attr"`
	ConnectedAP []*ConnectedAP `xml:"ConnectedAP"`
}

type ConnectedAP struct {
	IedName  string   `xml:"iedName,attr"`
	APName   string   `xml:"apName,attr"`
	Address  *Address `xml:"Address"`
	GESNodes []*GSE   `xml:"GSE"`
	SMVNodes []*SMV   `xml:"SMV"`
}

type Address struct {
	AddressParameters []*AddressParameter `xml:"P"`
}

type GSE struct {
	LdInst     string         `xml:"ldInst,attr"`
	CbName     string         `xml:"cbName,attr"`
	MinTimeVal *Val           `xml:"MinTime"`
	MaxTimeVal *Val           `xml:"MaxTime"`
	Address    *PhyComAddress `xml:"Address"`

	// custom
	MinTime int
	MaxTime int
}

type PhyComAddress struct {
	AddressParameters []*AddressParameter `xml:"P"`

	// custom
	VlanId       int64
	VlanPriority int
	AppId        int64
	MacAddress   []int
}

type SMV struct {
	LdInst  string         `xml:"ldInst,attr"`
	CbName  string         `xml:"cbName,attr"`
	Address *PhyComAddress `xml:"Address"`
}

type AddressParameter struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type IED struct {
	Name          string         `xml:"name,attr"`
	Type          string         `xml:"type,attr"`
	Manufacturer  string         `xml:"manufacturer,attr"`
	ConfigVersion string         `xml:"configVersion,attr"`
	Services      *Services      `xml:"Services"`
	AccessPoints  []*AccessPoint `xml:"AccessPoint"`
}

type Services struct {
	ConfDataSet       *ConfDataSet       `xml:"ConfDataSet"`
	ConfReportControl *ConfReportControl `xml:"ConfReportControl"`
	ReportSettings    *ReportSettings    `xml:"ReportSettings"`
	ConfLNs           *ConfLNs           `xml:"ConfLNs"`
}

type ConfDataSet struct {
	Max           int `xml:"max,attr"`
	MaxAttributes int `xml:"maxAttributes,attr"`
}

type ConfReportControl struct {
	Max int `xml:"max,attr"`
}

type ReportSettings struct {
	CbName    string `xml:"cbName,attr"`
	DatSet    string `xml:"datSet,attr"`
	RptID     string `xml:"rptID,attr"`
	OptFields string `xml:"optFields,attr"`
	BufTime   string `xml:"bufTime,attr"`
	TrgOps    string `xml:"trgOps,attr"`
	IntgPd    string `xml:"intgPd,attr"`
	Owner     bool   `xml:"owner,attr"`
}

type ConfLNs struct {
	FixPrefix bool `xml:"fixPrefix,attr"`
	FixLnInst bool `xml:"fixLnInst,attr"`
}

type AccessPoint struct {
	Name   string  `xml:"name,attr"`
	Server *Server `xml:"Server"`
}

type Server struct {
	Authentication *Authentication  `xml:"Authentication"`
	LogicalDevices []*LogicalDevice `xml:"LDevice" validate:"nonzero"`
}

type Authentication struct {
	None bool `xml:"none,attr"`
}

type LogicalDevice struct {
	Inst   string         `xml:"inst,attr"`
	LdName string         `xml:"ldName,attr"`
	LN0    *LogicalNode   `xml:"LN0" validate:"nonzero"`
	LNodes []*LogicalNode `xml:"LN"`

	// custom
	LogicalNodes []*LogicalNode // 将LN0、LN合并到一个数组中
}

type LogicalNode struct {
	Prefix  string `xml:"prefix,attr"`
	Inst    string `xml:"inst,attr"`
	LnClass string `xml:"lnClass,attr"`
	LnType  string `xml:"lnType,attr"`
	Desc    string `xml:"desc,attr"`

	DataSets                  []*DataSet             `xml:"DataSet"`
	ReportControlBlocks       []*ReportControl       `xml:"ReportControl"`
	GSEControlBlocks          []*GSEControl          `xml:"GSEControl"`
	SMVControlBlocks          []*SampledValueControl `xml:"SampledValueControl"`
	LogControlBlocks          []*LogControl          `xml:"LogControl"`
	Logs                      []*Log                 `xml:"Log"`
	SettingGroupControlBlocks []*SettingControl      `xml:"SettingControl"`
	DOINodes                  []*DOINode             `xml:"DOI"`

	// custom
	SclType     SclType
	DataObjects []*DataObject
}

type DataSet struct {
	Name string  `xml:"name,attr"`
	Desc string  `xml:"desc,attr"`
	FCDA []*FCDA `xml:"FCDA"`
}

type FCDA struct {
	LdInst  string `xml:"ldInst,attr"`
	Prefix  string `xml:"prefix,attr"`
	LnInst  string `xml:"lnInst,attr"`
	LnClass string `xml:"lnClass,attr"`
	DoName  string `xml:"doName,attr"`
	DaName  string `xml:"daName,attr"`
	Fc      string `xml:"fc,attr"`
}

type ReportControl struct {
	Name           string          `xml:"name,attr"`
	Desc           string          `xml:"desc,attr"`
	DatSet         string          `xml:"datSet,attr"`
	RptID          string          `xml:"rptID,attr"`
	ConfRev        string          `xml:"confRev,attr"`
	Buffered       bool            `xml:"buffered,attr"`
	BufTime        int             `xml:"bufTime,attr"`
	IntgPd         string          `xml:"intgPd,attr"`
	IndexedStr     string          `xml:"indexed,attr"`
	TriggerOptions *TriggerOptions `xml:"TrgOps"`
	OptionFields   *OptionFields   `xml:"OptFields"`
	RptEnabled     *RptEnabled     `xml:"RptEnabled"`

	// custom
	Indexed bool
}

type GSEControl struct {
	AppID     string `xml:"appID,attr"`
	Name      string `xml:"name,attr"`
	Desc      string `xml:"desc,attr"`
	Type      string `xml:"type,attr"`
	DatSet    string `xml:"datSet,attr"`
	ConfRev   int    `xml:"confRev,attr"`
	FixedOffs bool   `xml:"fixedOffs,attr"`
}

type SampledValueControl struct {
	Name      string   `xml:"name,attr"`
	Desc      string   `xml:"desc,attr"`
	DatSet    string   `xml:"datSet,attr"`
	SmvID     string   `xml:"smvID,attr"`
	SmpRate   int      `xml:"smpRate,attr"`
	NofASDU   int      `xml:"nofASDU,attr"`
	ConfRev   int      `xml:"confRev,attr"`
	SmpModStr string   `xml:"smpMod,attr"`
	Multicast bool     `xml:"multicast,attr"`
	SmvOpts   *SmvOpts `xml:"SmvOpts"`

	// custom
	SmpMod SmpMod
}

type SmvOpts struct {
	RefreshTime        bool `xml:"refreshTime,attr"`
	SampleSynchronized bool `xml:"sampleSynchronized,attr"`
	Security           bool `xml:"security,attr"`
	DataSet            bool `xml:"dataSet,attr"`
	SampleRate         bool `xml:"sampleRate,attr"`
}

type LogControl struct {
	Name           string          `xml:"name,attr"`
	Desc           string          `xml:"desc,attr"`
	DatSet         string          `xml:"datSet,attr"`
	LdInst         string          `xml:"ldInst,attr"`
	Prefix         string          `xml:"prefix,attr"`
	LnClass        string          `xml:"lnClass,attr"`
	LnInst         string          `xml:"lnInst,attr"`
	IntgPd         int             `xml:"intgPd,attr"`
	ReasonCode     bool            `xml:"reasonCode,attr"`
	LogName        string          `xml:"logName,attr"`
	LogEna         bool            `xml:"logEna,attr"`
	TriggerOptions *TriggerOptions `xml:"TrgOps"`
}

type Log struct {
	Name string `xml:"name,attr"`
}

type TriggerOptions struct {
	Dchg   bool `xml:"dchg,attr"`
	Qchg   bool `xml:"qchg,attr"`
	Dupd   bool `xml:"dupd,attr"`
	Period bool `xml:"period,attr"`
	Gi     bool `xml:"gi,attr"`
}

type OptionFields struct {
	SeqNum     bool `xml:"seqNum,attr"`
	TimeStamp  bool `xml:"timeStamp,attr"`
	DataSet    bool `xml:"dataSet,attr"`
	ReasonCode bool `xml:"reasonCode,attr"`
	DataRef    bool `xml:"dataRef,attr"`
	EntryID    bool `xml:"entryID,attr"`
	ConfigRef  bool `xml:"configRef,attr"`
	BufOvfl    bool `xml:"bufOvfl,attr"`
}

type RptEnabled struct {
	MaxStr    string      `xml:"max,attr"`
	Desc      string      `xml:"desc,attr"`
	ClientLNs []*ClientLN `xml:"ClientLN"`

	// custom
	Max int
}

type ClientLN struct {
	IedName string `xml:"iedName,attr"`
	ApRef   string `xml:"apRef,attr"`
	LdInst  string `xml:"ldInst,attr"`
	Prefix  string `xml:"prefix,attr"`
	LnClass string `xml:"lnClass,attr"`
	InInst  string `xml:"lnInst,attr"`
	Desc    string `xml:"desc,attr"`
}

type SettingControl struct {
	ActSG    int    `xml:"actSG,attr"`
	NumOfSGs int    `xml:"numOfSGs,attr"`
	Desc     string `xml:"desc,attr"`
}

type DOINode struct {
	Name  string `xml:"name,attr"`
	Desc  string `xml:"desc,attr"`
	SAddr string `xml:"sAddr,attr"`
	Val   *Val   `xml:"Val"`

	SDINodes []*DOINode `xml:"SDI"`
	DAINodes []*DOINode `xml:"DAI"`
}

type DataTypeTemplates struct {
	LogicalNodeTypes   []*LogicalNodeType   `xml:"LNodeType"`
	DataObjectTypes    []*DataObjectType    `xml:"DOType"`
	DataAttributeTypes []*DataAttributeType `xml:"DAType"`
	EnumTypes          []*EnumerationType   `xml:"EnumType"`
	TypeDeclarations   []SclType            // 将LNodeType、DOType、DAType、EnumType合并到一个数组中
}

type sclType struct {
	Id          string `xml:"id,attr"`
	Description string `xml:"desc,attr"`

	// custom
	Used bool
}

type LogicalNodeType struct {
	sclType
	LnClass               string                  `xml:"lnClass,attr"`
	DataObjectDefinitions []*DataObjectDefinition `xml:"DO"`
}

type DataObjectType struct {
	sclType
	Cdc            string                     `xml:"cdc,attr"`
	DataAttributes []*DataAttributeDefinition `xml:"DA"`
	SubDataObjects []*DataObjectDefinition    `xml:"SDO"`
}

type DataAttributeType struct {
	sclType
	SubDataAttributes []*DataAttributeDefinition `xml:"BDA"`
}

type EnumerationType struct {
	sclType
	EnumValues []*EnumerationValue `xml:"EnumVal"`
}

type EnumerationValue struct {
	Ord          int    `xml:"ord,attr"`
	SymbolicName string `xml:",chardata"`
}

type DataObjectDefinition struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	Transient bool   `xml:"transient,attr"`
	Count     int    `xml:"count,attr"`

	// custom
	SclType SclType
}

type DataAttributeDefinition struct {
	Name        string `xml:"name,attr"`
	Fc          string `xml:"fc,attr"`
	Type        string `xml:"type,attr"`
	BType       string `xml:"bType,attr"`
	Count       int    `xml:"count,attr"`
	DchgTrigger bool   `xml:"dchg,attr"`
	DupdTrigger bool   `xml:"dupd,attr"`
	QchgTrigger bool   `xml:"qchg,attr"`
	Val         *Val   `xml:"Val"`

	// custom
	AttributeType  AttributeType
	TriggerOptions *TriggerOptions
	Value          *DataModelValue
}

type Val struct {
	Value string `xml:",chardata"`
}

func (s *sclType) GetId() string {
	return s.Id
}

func (s *sclType) GetDesc() string {
	return s.Description
}

func (s *sclType) GetUsed() bool {
	return s.Used
}

func (s *sclType) SetUsed(used bool) {
	s.Used = used
}

func (s *SCL) getFirstIed() *IED {
	return s.IEDs[0]
}

func (s *SCL) getIedByName(iedName string) *IED {
	for _, ied := range s.IEDs {
		if ied.Name == iedName {
			return ied
		}
	}

	return nil
}

func (i *IED) getFirstAccessPoint() *AccessPoint {
	return i.AccessPoints[0]
}

func (i *IED) getAccessPointByName(apName string) *AccessPoint {
	for _, ap := range i.AccessPoints {
		if ap.Name == apName {
			return ap
		}
	}

	return nil
}

func (d *DataTypeTemplates) lookupType(id string) SclType {
	for _, declaration := range d.TypeDeclarations {
		if declaration.GetId() == id {
			return declaration
		}
	}

	return nil
}

func (e *EnumerationType) getOrdByEnumString(enumString string) (int, error) {
	for _, item := range e.EnumValues {
		if item.SymbolicName == enumString {
			return item.Ord, nil
		}
	}

	return 0, fmt.Errorf("enum has no value %s", enumString)
}

func (e *EnumerationType) isValidOrdValue(ordValue int) bool {
	for _, item := range e.EnumValues {
		if item.Ord == ordValue {
			return true
		}
	}

	return false
}

func (ln *LogicalNode) GetName() string {
	var name string

	if ln.Prefix != "" {
		name += ln.Prefix
	}

	name += ln.LnClass
	name += ln.Inst

	return name
}

func (ln *LogicalNode) GetSclType() SclType {
	return ln.SclType
}

func (ln *LogicalNode) GetChildByName(childName string) DataModelNode {
	if ln.DataObjects != nil {
		for _, do := range ln.DataObjects {
			if do.GetName() == childName {
				return do
			}
		}
	}

	return nil
}

func (that *Communication) getConnectedAP(apName string) *ConnectedAP {
	if that.SubNetworks != nil {

		for _, subNetwork := range that.SubNetworks {
			if subNetwork.ConnectedAP != nil {

				for _, connectedAP := range subNetwork.ConnectedAP {
					if connectedAP.APName == apName {
						return connectedAP
					}
				}
			}
		}
	}

	return nil
}

func (that *Communication) getIpAddressByIedName(iedName, apRef string) string {
	if that.SubNetworks != nil {

		for _, subNetwork := range that.SubNetworks {
			if subNetwork.ConnectedAP != nil {

				for _, ap := range subNetwork.ConnectedAP {
					if apRef != "" {
						isMatching := false

						if ap.APName == apRef {
							isMatching = true
						}

						if !isMatching {
							continue
						}
					}

					if ap.IedName != "" && ap.IedName == iedName {
						if ap.Address != nil && ap.Address.AddressParameters != nil {
							for _, p := range ap.Address.AddressParameters {
								if p.Type == "IP" {
									return p.Value
								}
							}
						}
					}
				}

			}
		}
	}

	return ""
}

func (that *TriggerOptions) GetIntValue() int {
	intValue := 0

	if that.Dchg {
		intValue += 1
	}
	if that.Qchg {
		intValue += 2
	}
	if that.Dupd {
		intValue += 4
	}
	if that.Period {
		intValue += 8
	}
	if that.Gi {
		intValue += 16
	}

	return intValue
}

func (that *ConnectedAP) LookupGSE(logicalDeviceName, name string) *GSE {
	for _, gse := range that.GESNodes {
		if gse.LdInst == logicalDeviceName && gse.CbName == name {
			return gse
		}
	}
	return nil
}

func (that *ConnectedAP) LookupSMV(logicalDeviceName, name string) *SMV {
	for _, smv := range that.SMVNodes {
		if smv.LdInst == logicalDeviceName && smv.CbName == name {
			return smv
		}
	}
	return nil
}

func (that *SmvOpts) GetIntValue() int {
	intValue := 0

	if that.RefreshTime {
		intValue += 1
	}
	if that.SampleSynchronized {
		intValue += 2
	}
	if that.SampleRate {
		intValue += 4
	}
	if that.DataSet {
		intValue += 8
	}
	if that.Security {
		intValue += 16
	}

	return intValue
}
