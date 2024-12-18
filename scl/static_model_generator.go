package scl

import (
	"fmt"
	"github.com/spf13/cast"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

type StaticModelGenerator struct {
	_scl           *SCL
	iedName        string
	apName         string
	outFileName    string
	outDir         string
	modelPrefix    string
	initializeOnce bool

	hDefineName string
	cOut        *PrintStream
	hOut        *PrintStream

	initializerBuffer *strings.Builder

	reportControlBlocks      *strings.Builder
	rcbVariableNames         []string
	currentRcbVariableNumber int

	gseControlBlocks         *strings.Builder
	gseVariableNames         []string
	currentGseVariableNumber int

	smvControlBlocks          *strings.Builder
	smvVariableNames          []string
	currentSvCBVariableNumber int

	settingGroupControlBlocks *strings.Builder
	sgcbVariableNames         []string
	currentSGCBVariableNumber int

	logControlBlocks         *strings.Builder
	lcbVariableNames         []string
	currentLcbVariableNumber int

	logs                     *strings.Builder
	logVariableNames         []string
	currentLogVariableNumber int

	ied         *IED
	hasOwner    bool
	accessPoint *AccessPoint
	connectedAP *ConnectedAP

	variablesList []string

	dataSetNames []string
}

func NewStaticModelGenerator(scl *SCL, iedName, ap, outDir, outFileName, modelPrefix string, initializeOnce bool) *StaticModelGenerator {
	hDefineName := strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(outFileName), ".", "_"), "-", "_") + "_H_"
	if strings.LastIndex(hDefineName, "/") >= 0 {
		hDefineName = hDefineName[strings.LastIndex(hDefineName, "/")+1:]
	}

	return &StaticModelGenerator{
		_scl:                      scl,
		iedName:                   iedName,
		apName:                    ap,
		outFileName:               outFileName,
		outDir:                    outDir,
		modelPrefix:               modelPrefix,
		initializeOnce:            initializeOnce,
		hDefineName:               hDefineName,
		cOut:                      &PrintStream{&strings.Builder{}},
		hOut:                      &PrintStream{&strings.Builder{}},
		initializerBuffer:         &strings.Builder{},
		reportControlBlocks:       &strings.Builder{},
		rcbVariableNames:          make([]string, 0),
		gseControlBlocks:          &strings.Builder{},
		gseVariableNames:          make([]string, 0),
		smvControlBlocks:          &strings.Builder{},
		smvVariableNames:          make([]string, 0),
		settingGroupControlBlocks: &strings.Builder{},
		sgcbVariableNames:         make([]string, 0),
		logControlBlocks:          &strings.Builder{},
		lcbVariableNames:          make([]string, 0),
		logs:                      &strings.Builder{},
		logVariableNames:          make([]string, 0),
		variablesList:             make([]string, 0),
		dataSetNames:              make([]string, 0),
	}
}

func (s *StaticModelGenerator) Generate() error {

	if s.iedName == "" {
		s.ied = s._scl.getFirstIed()
	} else {
		s.ied = s._scl.getIedByName(s.iedName)
	}

	if s.ied == nil {
		return fmt.Errorf("IED model not found in SCL file")
	}

	if s.ied.Services != nil && s.ied.Services.ReportSettings != nil {
		s.hasOwner = s.ied.Services.ReportSettings.Owner
	}

	if s.apName == "" {
		s.accessPoint = s.ied.getFirstAccessPoint()
	} else {
		s.accessPoint = s.ied.getAccessPointByName(s.apName)
	}

	if s._scl.Communication != nil {
		s.connectedAP = s._scl.Communication.getConnectedAP(s.accessPoint.Name)
	}

	s.printCFileHeader()
	s.printHeaderFileHeader()
	s.printForwardDeclarations(s.accessPoint.Server)
	if err := s.printDeviceModelDefinitions(); err != nil {
		return err
	}
	s.printInitializerFunction()
	s.printVariablePointerDefines()
	s.printHeaderFileFooter()

	// output to files
	if err := s.cOut.writeFile(filepath.Join(s.outDir, s.outFileName+".c")); err != nil {
		return err
	}
	if err := s.hOut.writeFile(filepath.Join(s.outDir, s.outFileName+".h")); err != nil {
		return err
	}
	return nil
}

func (s *StaticModelGenerator) printCFileHeader() {
	include := s.outFileName + ".h"
	if strings.LastIndex(include, "/") >= 0 {
		include = include[strings.LastIndex(include, "/")+1:]
	}

	s.cOut.println("/*")
	s.cOut.println(" * %s.c", s.outFileName)
	s.cOut.println(" *")
	s.cOut.println(" * automatically generated from %s", s._scl.FileFullPath)
	s.cOut.println(" */")
	s.cOut.println("#include \"" + include + "\"")
	s.cOut.printlnNone()
}

func (s *StaticModelGenerator) printHeaderFileHeader() {
	s.hOut.println("/*")
	s.hOut.println(" * %s.h", s.outFileName)
	s.hOut.println(" *")
	s.hOut.println(" * automatically generated from %s", s._scl.FileFullPath)
	s.hOut.println(" */\n")
	s.hOut.println("#ifndef %s", s.hDefineName)
	s.hOut.println("#define %s\n", s.hDefineName)
	s.hOut.println("#include <stdlib.h>")
	s.hOut.println("#include \"iec61850_model.h\"")
	s.hOut.printlnNone()
}

func (s *StaticModelGenerator) printForwardDeclarations(server *Server) {
	s.cOut.println("static void initializeValues();")
	s.hOut.println("extern IedModel %s;", s.modelPrefix)

	for _, lDevice := range server.LogicalDevices {
		ldName := s.modelPrefix + "_" + lDevice.Inst
		s.hOut.println("extern LogicalDevice %s;", ldName)

		for _, logicalNode := range lDevice.LogicalNodes {
			lnName := ldName + "_" + logicalNode.GetName()

			s.hOut.println("extern LogicalNode   %s;", lnName)
			s.printDataObjectForwardDeclarationsByDO(lnName, logicalNode.DataObjects, make(map[string]uint8))
		}
	}
}

func (s *StaticModelGenerator) printDataObjectForwardDeclarationsByDO(prefix string, dataObjects []*DataObject, doNameMap map[string]uint8) {
	for _, do := range dataObjects {
		doName := prefix + "_" + do.GetName()

		if _, exist := doNameMap[doName]; exist {
			fmt.Printf("extern DataObject %s already exists\n", doName)
			continue
		}

		doNameMap[doName] = 1
		s.hOut.println("extern DataObject    %s;", doName)

		if do.SubDataObjects != nil {
			s.printDataObjectForwardDeclarationsByDO(doName, do.SubDataObjects, doNameMap)
		}

		s.printDataObjectForwardDeclarationsByDA(doName, do.DataAttributes, make(map[string]uint8))
	}
}

func (s *StaticModelGenerator) printDataObjectForwardDeclarationsByDA(doName string, dataAttributes []*DataAttribute, daNameMap map[string]uint8) {
	for _, da := range dataAttributes {
		daName := doName + "_" + da.GetName()

		if da.FC == "SE" {
			if !strings.HasPrefix(daName, s.modelPrefix+"_SE_") {
				daName = daName[:9] + "SE_" + daName[9:]
			}
		}

		if _, exist := daNameMap[daName]; exist {
			fmt.Printf("extern DataAttribute %s already exists\n", daName)
			continue
		}

		daNameMap[daName] = 1
		s.hOut.println("extern DataAttribute %s;", daName)
		if da.SubDataAttributes != nil {
			s.printDataObjectForwardDeclarationsByDA(daName, da.SubDataAttributes, daNameMap)
		}
	}
}

func (s *StaticModelGenerator) printDeviceModelDefinitions() error {
	if err := s.printDataSets(); err != nil {
		return err
	}

	s.createLNSubVariableList(s.accessPoint.Server.LogicalDevices)
	logicalDevices := s.accessPoint.Server.LogicalDevices
	for i, logicalDevice := range logicalDevices {
		ldName := s.modelPrefix + "_" + logicalDevice.Inst
		s.variablesList = append(s.variablesList, ldName)

		s.cOut.println("\nLogicalDevice %s = {", ldName)
		s.cOut.println("    LogicalDeviceModelType,")
		s.cOut.println("    \"%s\",", logicalDevice.Inst)
		s.cOut.println("    (ModelNode*) &%s,", s.modelPrefix)

		if i < (len(logicalDevices) - 1) {
			s.cOut.println("    (ModelNode*) &%s_%s,", s.modelPrefix, logicalDevices[i+1].Inst)
		} else {
			s.cOut.println("    NULL,")
		}

		firstChildName := fmt.Sprintf("%s_%s", ldName, logicalDevice.LogicalNodes[0].GetName())
		s.cOut.println("    (ModelNode*) &%s", firstChildName)
		s.cOut.println("};\n")

		s.printLogicalNodeDefinitions(ldName, logicalDevice, logicalDevice.LogicalNodes)
	}

	for _, rcb := range s.rcbVariableNames {
		s.cOut.println("extern ReportControlBlock %s;", rcb)
	}
	s.cOut.printlnNone()
	s.cOut.println(s.reportControlBlocks.String())

	for _, smv := range s.smvVariableNames {
		s.cOut.println("extern SVControlBlock %s;", smv)
	}
	s.cOut.println(s.smvControlBlocks.String())

	for _, gcb := range s.gseVariableNames {
		s.cOut.println("extern GSEControlBlock %s;", gcb)
	}
	s.cOut.println(s.gseControlBlocks.String())

	for _, sgcb := range s.sgcbVariableNames {
		s.cOut.println("extern SettingGroupControlBlock %s;", sgcb)
	}
	s.cOut.println(s.settingGroupControlBlocks.String())

	for _, lcb := range s.lcbVariableNames {
		s.cOut.println("extern LogControlBlock %s;", lcb)
	}
	s.cOut.println(s.logControlBlocks.String())

	for _, log := range s.logVariableNames {
		s.cOut.println("extern Log %s;", log)
	}
	s.cOut.println(s.logs.String())

	firstLogicalDeviceName := logicalDevices[0].Inst
	s.cOut.println("\nIedModel %s = {", s.modelPrefix)
	s.cOut.println("    \"%s\",", s.ied.Name)
	s.cOut.println("    &%s_%s,", s.modelPrefix, firstLogicalDeviceName)

	if len(s.dataSetNames) > 0 {
		s.cOut.println("    &%s,", s.dataSetNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.rcbVariableNames) > 0 {
		s.cOut.println("    &%s,", s.rcbVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.gseVariableNames) > 0 {
		s.cOut.println("    &%s,", s.gseVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.smvVariableNames) > 0 {
		s.cOut.println("    &%s,", s.smvVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.sgcbVariableNames) > 0 {
		s.cOut.println("    &%s,", s.sgcbVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.lcbVariableNames) > 0 {
		s.cOut.println("    &%s,", s.lcbVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	if len(s.logVariableNames) > 0 {
		s.cOut.println("    &%s,", s.logVariableNames[0])
	} else {
		s.cOut.println("    NULL,")
	}

	s.cOut.println("    initializeValues};")
	return nil
}

func (s *StaticModelGenerator) printInitializerFunction() {
	s.cOut.println("\nstatic void")
	s.cOut.println("initializeValues()")
	s.cOut.println("{")
	s.cOut.print(s.initializerBuffer.String())
	s.cOut.println("}")
}

func (s *StaticModelGenerator) printVariablePointerDefines() {
	s.hOut.println("\n\n")
	for _, variableName := range s.variablesList {
		name := strings.ToUpper(s.modelPrefix) + variableName[len(s.modelPrefix):]
		s.hOut.println(fmt.Sprintf("#define %s (&%s)", name, variableName))
	}
}

func (s *StaticModelGenerator) printHeaderFileFooter() {
	s.hOut.printlnNone()
	s.hOut.println("#endif /* %s */\n", s.hDefineName)
}

func (s *StaticModelGenerator) printLogicalNodeDefinitions(ldName string, logicalDevice *LogicalDevice, logicalNodes []*LogicalNode) {
	for i, logicalNode := range logicalNodes {
		lnName := ldName + "_" + logicalNode.GetName()

		s.variablesList = append(s.variablesList, lnName)
		// Print LogicalNode definition
		s.cOut.println("LogicalNode %s = {", lnName)
		s.cOut.println("    LogicalNodeModelType,")
		s.cOut.println("    \"%s\",", logicalNode.GetName())
		s.cOut.println("    (ModelNode*) &%s,", ldName)

		if i < len(logicalNodes)-1 {
			nextNodeName := fmt.Sprintf("%s_%s", ldName, logicalNodes[i+1].GetName())
			s.cOut.println("    (ModelNode*) &%s,", nextNodeName)
		} else {
			s.cOut.println("    NULL,")
		}

		// First child data object
		firstChildName := fmt.Sprintf("%s_%s", lnName, logicalNode.DataObjects[0].GetName())
		s.cOut.println("    (ModelNode*) &%s,", firstChildName)
		s.cOut.println("};\n")

		s.printDataObjectDefinitions(lnName, logicalNode.DataObjects, "", false)
		s.printReportControlBlocks(lnName, logicalNode)
		s.printLogControlBlocks(lnName, logicalNode, logicalDevice)
		s.printLogs(lnName, logicalNode)
		s.printGSEControlBlocks(ldName, lnName, logicalNode)
		s.printSVControlBlocks(ldName, lnName, logicalNode)
		s.printSettingControlBlock(lnName, logicalNode)
	}
}

func (s *StaticModelGenerator) printSettingControlBlock(lnPrefix string, logicalNode *LogicalNode) {
	settingControls := logicalNode.SettingGroupControlBlocks

	if len(settingControls) > 0 {
		sgcb := settingControls[0]
		sgcbVariableName := lnPrefix + "_sgcb"

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("\nSettingGroupControlBlock %s = {", sgcbVariableName))
		sb.WriteString(fmt.Sprintf("&%s, ", lnPrefix))
		sb.WriteString(fmt.Sprintf("%d, %d, 0, false, 0, 0, ", sgcb.ActSG, sgcb.NumOfSGs))

		if s.currentSGCBVariableNumber < len(s.sgcbVariableNames)-1 {
			sb.WriteString(fmt.Sprintf("&%s", s.sgcbVariableNames[s.currentSGCBVariableNumber+1]))
		} else {
			sb.WriteString("NULL")
		}

		sb.WriteString("};\n")
		s.settingGroupControlBlocks.WriteString(sb.String())

		s.currentSGCBVariableNumber++
	}
}

func (s *StaticModelGenerator) printSVControlBlocks(ldName, lnPrefix string, logicalNode *LogicalNode) {
	svControlBlocks := logicalNode.SMVControlBlocks

	// Strip "iedModel_" from ldName
	ldNameComponents := strings.Split(ldName, "_")
	logicalDeviceName := ldNameComponents[1]

	smvControlNumber := 0

	for _, svCB := range svControlBlocks {

		var svAddress *PhyComAddress
		smv := s.connectedAP.LookupSMV(logicalDeviceName, svCB.Name) // Assume lookupSMV is implemented
		if smv != nil {
			svAddress = smv.Address
		}

		var svStringBuilder strings.Builder
		var phyComAddrName string

		if svAddress != nil {
			phyComAddrName = fmt.Sprintf("%s_smv%d_address", lnPrefix, smvControlNumber)

			svStringBuilder.WriteString(fmt.Sprintf("\nstatic PhyComAddress %s = {\n", phyComAddrName))
			svStringBuilder.WriteString(fmt.Sprintf("  %d,\n", svAddress.VlanPriority))
			svStringBuilder.WriteString(fmt.Sprintf("  %d,\n", svAddress.VlanId))
			svStringBuilder.WriteString(fmt.Sprintf("  %d,\n", svAddress.AppId))
			svStringBuilder.WriteString("  {")
			for i, mac := range svAddress.MacAddress {
				svStringBuilder.WriteString(fmt.Sprintf("0x%s", strconv.FormatInt(int64(mac), 16)))
				if i < len(svAddress.MacAddress)-1 {
					svStringBuilder.WriteString(", ")
				} else {
					svStringBuilder.WriteString("}\n")
				}
			}
			svStringBuilder.WriteString("};\n\n")
		}

		smvVariableName := fmt.Sprintf("%s_smv%d", lnPrefix, smvControlNumber)

		svStringBuilder.WriteString(fmt.Sprintf("SVControlBlock %s = {", smvVariableName))
		svStringBuilder.WriteString(fmt.Sprintf("&%s, ", lnPrefix))
		svStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", svCB.Name))

		if svCB.SmvID == "" {
			svStringBuilder.WriteString("NULL, ")
		} else {
			svStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", svCB.SmvID))
		}

		if svCB.DatSet != "" {
			svStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", svCB.DatSet))
		} else {
			svStringBuilder.WriteString("NULL, ")
		}

		svStringBuilder.WriteString(fmt.Sprintf("%d, ", svCB.SmvOpts.GetIntValue()))
		svStringBuilder.WriteString(fmt.Sprintf("%d, ", svCB.SmpMod))
		svStringBuilder.WriteString(fmt.Sprintf("%d, ", svCB.SmpRate))
		svStringBuilder.WriteString(fmt.Sprintf("%d, ", svCB.ConfRev))

		if svAddress != nil {
			svStringBuilder.WriteString(fmt.Sprintf("&%s, ", phyComAddrName))
		} else {
			svStringBuilder.WriteString("NULL, ")
		}

		if svCB.Multicast {
			svStringBuilder.WriteString("false, ")
		} else {
			svStringBuilder.WriteString("true, ")
		}

		svStringBuilder.WriteString(fmt.Sprintf("%d, ", svCB.NofASDU))
		s.currentSvCBVariableNumber++

		if s.currentSvCBVariableNumber < len(s.smvVariableNames) {
			svStringBuilder.WriteString(fmt.Sprintf("&%s", s.smvVariableNames[s.currentSvCBVariableNumber]))
		} else {
			svStringBuilder.WriteString("NULL")
		}

		svStringBuilder.WriteString("};\n")
		s.smvControlBlocks.WriteString(svStringBuilder.String())

		smvControlNumber++
	}
}

func (s *StaticModelGenerator) printGSEControlBlocks(ldName, lnPrefix string, logicalNode *LogicalNode) {
	gseControlBlocks := logicalNode.GSEControlBlocks

	// Strip "iedModel_" from ldName
	ldNameComponents := strings.Split(ldName, "_")
	logicalDeviceName := ldNameComponents[1]

	gseControlNumber := 0

	for _, gseControlBlock := range gseControlBlocks {
		gse := s.connectedAP.LookupGSE(logicalDeviceName, gseControlBlock.Name) // Assume lookupGSE is implemented
		if gse != nil {
			var gseStringBuilder strings.Builder

			var phyComAddrName string
			if gse.Address != nil {
				phyComAddrName = fmt.Sprintf("%s_gse%d_address", lnPrefix, gseControlNumber)

				gseStringBuilder.WriteString(fmt.Sprintf("\nstatic PhyComAddress %s = {\n", phyComAddrName))
				gseStringBuilder.WriteString(fmt.Sprintf("  %d,\n", gse.Address.VlanPriority))
				gseStringBuilder.WriteString(fmt.Sprintf("  %d,\n", gse.Address.VlanId))
				gseStringBuilder.WriteString(fmt.Sprintf("  %d,\n", gse.Address.AppId))
				gseStringBuilder.WriteString("  {")
				for i, mac := range gse.Address.MacAddress {
					gseStringBuilder.WriteString(fmt.Sprintf("0x%s", strconv.FormatInt(int64(mac), 16)))
					if i < len(gse.Address.MacAddress)-1 {
						gseStringBuilder.WriteString(", ")
					}
				}
				gseStringBuilder.WriteString("}\n};\n\n")
			}

			gseVariableName := fmt.Sprintf("%s_gse%d", lnPrefix, gseControlNumber)

			gseStringBuilder.WriteString(fmt.Sprintf("GSEControlBlock %s = {", gseVariableName))
			gseStringBuilder.WriteString(fmt.Sprintf("&%s, ", lnPrefix))
			gseStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", gseControlBlock.Name))

			if gseControlBlock.AppID == "" {
				gseStringBuilder.WriteString("NULL, ")
			} else {
				gseStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", gseControlBlock.AppID))
			}

			if gseControlBlock.DatSet != "" {
				gseStringBuilder.WriteString(fmt.Sprintf("\"%s\", ", gseControlBlock.DatSet))
			} else {
				gseStringBuilder.WriteString("NULL, ")
			}

			gseStringBuilder.WriteString(fmt.Sprintf("%d, ", gseControlBlock.ConfRev))
			gseStringBuilder.WriteString(fmt.Sprintf("%t, ", gseControlBlock.FixedOffs))

			if gse.Address != nil {
				gseStringBuilder.WriteString(fmt.Sprintf("&%s, ", phyComAddrName))
			} else {
				gseStringBuilder.WriteString("NULL, ")
			}

			gseStringBuilder.WriteString(fmt.Sprintf("%d, ", gse.MinTime))
			gseStringBuilder.WriteString(fmt.Sprintf("%d, ", gse.MaxTime))

			s.currentGseVariableNumber++

			if s.currentGseVariableNumber < len(s.gseVariableNames) {
				gseStringBuilder.WriteString(fmt.Sprintf("&%s", s.gseVariableNames[s.currentGseVariableNumber]))
			} else {
				gseStringBuilder.WriteString("NULL")
			}

			gseStringBuilder.WriteString("};\n")
			s.gseControlBlocks.WriteString(gseStringBuilder.String())
			gseControlNumber++
		} else {
			fmt.Printf("GSE not found for GoCB %s\n", gseControlBlock.Name)
		}
	}
}

func (s *StaticModelGenerator) printLogs(lnPrefix string, logicalNode *LogicalNode) {
	for logNumber, log := range logicalNode.Logs {
		s.printLog(lnPrefix, log, logNumber)
	}
}

func (s *StaticModelGenerator) printLog(lnPrefix string, log *Log, logNumber int) {
	logVariableName := fmt.Sprintf("%s_log%d", lnPrefix, logNumber)

	var logString strings.Builder
	logString.WriteString(fmt.Sprintf("Log %s = {", logVariableName))

	// Add Logical Node reference
	logString.WriteString(fmt.Sprintf("&%s, ", lnPrefix))

	// Add Log name
	logString.WriteString(fmt.Sprintf("\"%s\", ", log.Name))

	// Add reference to the next log
	s.currentLogVariableNumber++
	if s.currentLogVariableNumber < len(s.logVariableNames) {
		logString.WriteString(fmt.Sprintf("&%s", s.logVariableNames[s.currentLogVariableNumber]))
	} else {
		logString.WriteString("NULL")
	}

	logString.WriteString("};\n")

	// Append to logs buffer
	s.logs.WriteString(logString.String())
}

func (s *StaticModelGenerator) printLogControlBlocks(lnPrefix string, logicalNode *LogicalNode, logicalDevice *LogicalDevice) {
	for lcbNumber, lcb := range logicalNode.LogControlBlocks {
		s.printLogControlBlock(logicalDevice, lnPrefix, lcb, lcbNumber)
	}
}

func (s *StaticModelGenerator) printLogControlBlock(logicalDevice *LogicalDevice, lnPrefix string, lcb *LogControl, lcbNumber int) {
	lcbVariableName := fmt.Sprintf("%s_lcb%d", lnPrefix, lcbNumber)

	var lcbString strings.Builder
	lcbString.WriteString(fmt.Sprintf("LogControlBlock %s = {", lcbVariableName))

	// Add Logical Node reference
	lcbString.WriteString(fmt.Sprintf("&%s, ", lnPrefix))

	// Add LogControl name
	lcbString.WriteString(fmt.Sprintf("\"%s\", ", lcb.Name))

	// Add DataSet
	if lcb.DatSet == "" {
		lcbString.WriteString("NULL, ")
	} else {
		lcbString.WriteString(fmt.Sprintf("\"%s\", ", lcb.DatSet))
	}

	// Build logRef
	var logRef string
	if lcb.LdInst == "" {
		logRef = logicalDevice.Inst + "/"
	} else {
		logRef = lcb.LdInst + "/"
	}

	if lcb.LnClass == "LLN0" {
		logRef += "LLN0$"
	} else {
		logRef += lcb.LnClass + "$"
	}

	// Add logRef and logName
	if lcb.LogName != "" {
		lcbString.WriteString(fmt.Sprintf("\"%s%s\", ", logRef, lcb.LogName))
	} else {
		lcbString.WriteString("NULL, ")
	}

	// Add TriggerOptions
	triggerOps := 0
	if lcb.TriggerOptions != nil {
		triggerOps = lcb.TriggerOptions.GetIntValue()
	}
	if triggerOps >= 16 {
		triggerOps -= 16
	}
	lcbString.WriteString(fmt.Sprintf("%d, ", triggerOps))

	// Add Integration Period
	if lcb.IntgPd != 0 {
		lcbString.WriteString(fmt.Sprintf("%d, ", lcb.IntgPd))
	} else {
		lcbString.WriteString("0, ")
	}

	// Add Log Enabled
	if lcb.LogEna {
		lcbString.WriteString("true, ")
	} else {
		lcbString.WriteString("false, ")
	}

	// Add Reason Code
	if lcb.ReasonCode {
		lcbString.WriteString("true, ")
	} else {
		lcbString.WriteString("false, ")
	}

	// Add reference to next LogControlBlock
	s.currentLcbVariableNumber++
	if s.currentLcbVariableNumber < len(s.lcbVariableNames) {
		lcbString.WriteString(fmt.Sprintf("&%s", s.lcbVariableNames[s.currentLcbVariableNumber]))
	} else {
		lcbString.WriteString("NULL")
	}

	lcbString.WriteString("};\n")

	// Append to logControlBlocks
	s.logControlBlocks.WriteString(lcbString.String())
}

// printReportControlBlocks handles printing of report control blocks in the logical node
func (s *StaticModelGenerator) printReportControlBlocks(lnPrefix string, logicalNode *LogicalNode) {
	reportControlBlocks := logicalNode.ReportControlBlocks

	reportsCount := len(reportControlBlocks)
	reportNumber := 0

	for _, rcb := range reportControlBlocks {
		if rcb.Indexed {
			maxInstances := 1
			var clientLNs []*ClientLN

			if rcb.RptEnabled != nil {
				maxInstances = rcb.RptEnabled.Max
				clientLNs = rcb.RptEnabled.ClientLNs
			}

			for i := 0; i < maxInstances; i++ {
				index := fmt.Sprintf("%02d", i+1)

				clientAddress := make([]byte, 17)
				clientAddress[0] = 0

				if clientLNs != nil {
					if i < len(clientLNs) {
						clientLN := clientLNs[i]
						if clientLN != nil {
							iedName := clientLN.IedName
							apRef := clientLN.ApRef

							if iedName != "" {
								ipAddress := s._scl.Communication.getIpAddressByIedName(iedName, apRef)

								// Resolve IP Address (IPv4 or IPv6)
								inetAddr, err := net.ResolveIPAddr("ip", ipAddress)
								if err == nil {
									if inetAddr.IP.To4() != nil { // IPv4
										clientAddress[0] = 4
										copy(clientAddress[1:], inetAddr.IP.To4())
									} else { // IPv6
										clientAddress[0] = 6
										copy(clientAddress[1:], inetAddr.IP.To16())
									}
								} else {
									// Handle the error appropriately
									fmt.Println("Error resolving IP address:", err)
								}
							}
						}
					}
				}

				// Print the report control block instance
				s.printReportControlBlockInstance(lnPrefix, rcb, index, reportNumber, reportsCount, clientAddress)
				reportNumber++
			}
		} else {
			// Default case for non-indexed report control blocks
			clientAddress := make([]byte, 17)
			clientAddress[0] = 0

			// Print the report control block instance
			s.printReportControlBlockInstance(lnPrefix, rcb, "", reportNumber, reportsCount, clientAddress)
			reportNumber++
		}
	}
}

func (s *StaticModelGenerator) printReportControlBlockInstance(lnPrefix string, rcb *ReportControl, index string, reportNumber int, reportsCount int, clientIpAddr []byte) {
	rcbVariableName := fmt.Sprintf("%s_report%d", lnPrefix, reportNumber)

	var rcbString strings.Builder
	rcbString.WriteString(fmt.Sprintf("ReportControlBlock %s = {", rcbVariableName))

	// Add Logical Node reference
	rcbString.WriteString(fmt.Sprintf("&%s, ", lnPrefix))

	// Add Report Name
	rcbString.WriteString(fmt.Sprintf("\"%s%s\", ", rcb.Name, index))

	// Add RptID
	if rcb.RptID == "" {
		rcbString.WriteString("NULL, ")
	} else {
		rcbString.WriteString(fmt.Sprintf("\"%s\", ", rcb.RptID))
	}

	// Add Buffered
	rcbString.WriteString(fmt.Sprintf("%t, ", rcb.Buffered))

	// Add DataSet
	if rcb.DatSet == "" {
		rcbString.WriteString("NULL, ")
	} else {
		rcbString.WriteString(fmt.Sprintf("\"%s\", ", rcb.DatSet))
	}

	// Add ConfRef
	if rcb.ConfRev == "" {
		rcbString.WriteString("0, ")
	} else {
		rcbString.WriteString(fmt.Sprintf("%s, ", rcb.ConfRev))
	}

	// Add TriggerOptions
	triggerOps := 16
	if rcb.TriggerOptions != nil {
		triggerOps = rcb.TriggerOptions.GetIntValue()
	}
	if s.hasOwner {
		triggerOps += 64
	}
	rcbString.WriteString(fmt.Sprintf("%d, ", triggerOps))

	// Add OptionFields
	options := 0
	if rcb.OptionFields != nil {
		if rcb.OptionFields.SeqNum {
			options += 1
		}
		if rcb.OptionFields.TimeStamp {
			options += 2
		}
		if rcb.OptionFields.ReasonCode {
			options += 4
		}
		if rcb.OptionFields.DataSet {
			options += 8
		}
		if rcb.OptionFields.DataRef {
			options += 16
		}
		if rcb.OptionFields.BufOvfl {
			options += 32
		}
		if rcb.OptionFields.EntryID {
			options += 64
		}
		if rcb.OptionFields.ConfigRef {
			options += 128
		}
	} else {
		options = 32
	}
	rcbString.WriteString(fmt.Sprintf("%d, ", options))

	// Add BufferTime
	rcbString.WriteString(fmt.Sprintf("%d, ", rcb.BufTime))

	// Add IntegrityPeriod
	if rcb.IntgPd == "" {
		rcbString.WriteString("0, ")
	} else {
		rcbString.WriteString(fmt.Sprintf("%s, ", rcb.IntgPd))
	}

	// Add Client IP Address
	rcbString.WriteString("{")
	for i, byteVal := range clientIpAddr {
		rcbString.WriteString(fmt.Sprintf("0x%02X", byteVal))
		if i < len(clientIpAddr)-1 {
			rcbString.WriteString(", ")
		} else {
			rcbString.WriteString("}, ")
		}
	}

	// Add next ReportControlBlock reference
	s.currentRcbVariableNumber++
	if s.currentRcbVariableNumber < len(s.rcbVariableNames) {
		rcbString.WriteString(fmt.Sprintf("&%s", s.rcbVariableNames[s.currentRcbVariableNumber]))
	} else {
		rcbString.WriteString("NULL")
	}

	rcbString.WriteString("};\n")

	// Append to reportControlBlocks
	s.reportControlBlocks.WriteString(rcbString.String())
}

func (s *StaticModelGenerator) printDataObjectDefinitions(lnName string, dataObjects []*DataObject, dataAttributeSibling string, isTransient bool) {
	for i, dataObject := range dataObjects {
		doName := fmt.Sprintf("%s_%s", lnName, dataObject.Name)
		s.variablesList = append(s.variablesList, doName)

		// Print DataObject definition
		s.cOut.println("DataObject %s = {", doName)
		s.cOut.println("    DataObjectModelType,")
		s.cOut.println("    \"%s\",", dataObject.Name)
		s.cOut.println("    (ModelNode*) &%s,", lnName)

		// Determine sibling node
		if i < len(dataObjects)-1 {
			nextSibling := fmt.Sprintf("%s_%s", lnName, dataObjects[i+1].Name)
			s.cOut.println("    (ModelNode*) &%s,", nextSibling)
		} else if dataAttributeSibling != "" {
			s.cOut.println("    (ModelNode*) &%s,", dataAttributeSibling)
		} else {
			s.cOut.println("    NULL,")
		}

		// Determine first child node
		var firstSubDataObjectName, firstDataAttributeName string
		if len(dataObject.SubDataObjects) > 0 {
			firstSubDataObjectName = fmt.Sprintf("%s_%s", doName, dataObject.SubDataObjects[0].Name)
		}
		if len(dataObject.DataAttributes) > 0 {
			attribute := dataObject.DataAttributes[0]

			fDoName := doName
			if attribute.FC == "SE" {
				if !strings.HasPrefix(doName, s.modelPrefix+"_SE_") {
					fDoName = doName[:9] + "SE_" + doName[9:]
				}
			}
			firstDataAttributeName = fmt.Sprintf("%s_%s", fDoName, attribute.Name)
		}

		if firstSubDataObjectName != "" {
			s.cOut.println("    (ModelNode*) &%s,", firstSubDataObjectName)
		} else if firstDataAttributeName != "" {
			s.cOut.println("    (ModelNode*) &%s,", firstDataAttributeName)
		} else {
			s.cOut.println("    NULL,")
		}

		// Print count
		s.cOut.println("    %d", dataObject.Count)
		s.cOut.println("};\n")

		// Determine if the current DataObject is transient
		isDoTransient := isTransient || dataObject.Trans

		// Recursively print sub-data objects and attributes
		if len(dataObject.SubDataObjects) > 0 {
			s.printDataObjectDefinitions(doName, dataObject.SubDataObjects, firstDataAttributeName, isDoTransient)
		}
		if len(dataObject.DataAttributes) > 0 {
			s.printDataAttributeDefinitions(doName, dataObject.DataAttributes, isDoTransient)
		}
	}
}

func (s *StaticModelGenerator) printDataAttributeDefinitions(doName string, dataAttributes []*DataAttribute, isTransient bool) {
	for i, dataAttribute := range dataAttributes {
		daName := doName + "_" + dataAttribute.Name

		// Handle FunctionalConstraint "SE"
		if dataAttribute.FC == "SE" {
			if !strings.HasPrefix(daName, s.modelPrefix+"_SE_") {
				daName = daName[:9] + "SE_" + daName[9:]
			}
		}

		s.variablesList = append(s.variablesList, daName)

		// Print DataAttribute definition
		s.cOut.println("DataAttribute %s = {", daName)
		s.cOut.println("    DataAttributeModelType,")
		s.cOut.println("    \"%s\",", dataAttribute.Name)
		s.cOut.println("    (ModelNode*) &%s,", doName)

		// Sibling node
		if i < len(dataAttributes)-1 {
			sibling := dataAttributes[i+1]
			siblingDoName := doName
			if sibling.FC == "SE" {
				if !strings.HasPrefix(siblingDoName, s.modelPrefix+"_SE_") {
					siblingDoName = siblingDoName[:9] + "SE_" + siblingDoName[9:]
				}
			}
			s.cOut.println("    (ModelNode*) &%s_%s,", siblingDoName, sibling.Name)
		} else {
			s.cOut.println("    NULL,")
		}

		// First sub-data attribute
		if len(dataAttribute.SubDataAttributes) > 0 {
			s.cOut.println("    (ModelNode*) &%s_%s,", daName, dataAttribute.SubDataAttributes[0].Name)
		} else {
			s.cOut.println("    NULL,")
		}

		// Print Count, FunctionalConstraint, and Type
		s.cOut.println("    %d,", dataAttribute.Count)
		s.cOut.println("    IEC61850_FC_%s,", dataAttribute.FC)
		s.cOut.println("    IEC61850_%s,", dataAttribute.AttributeType.ToString())

		// Print trigger options
		s.cOut.print("    0")
		trgOps := dataAttribute.TriggerOptions
		if trgOps != nil {
			if trgOps.Dchg {
				s.cOut.print(" + TRG_OPT_DATA_CHANGED")
			}
			if trgOps.Dupd {
				s.cOut.print(" + TRG_OPT_DATA_UPDATE")
			}
			if trgOps.Qchg {
				s.cOut.print(" + TRG_OPT_QUALITY_CHANGED")
			}
		}
		if isTransient {
			s.cOut.print(" + TRG_OPT_TRANSIENT")
		}
		s.cOut.println(",")

		s.cOut.println("    NULL,")

		// Short address
		shortAddr := int64(0)
		if addr, err := cast.ToInt64E(dataAttribute.ShortAddress); err == nil {
			shortAddr = addr
		} else if dataAttribute.ShortAddress != "" {
			fmt.Printf("WARNING: short address \"%s\" is not valid for libIEC61850!\n", dataAttribute.ShortAddress)
		}
		s.cOut.println("    %d", shortAddr)
		s.cOut.println("};\n")

		// Recursive call for sub-data attributes
		if len(dataAttribute.SubDataAttributes) > 0 {
			s.printDataAttributeDefinitions(daName, dataAttribute.SubDataAttributes, isTransient)
		}

		// Process value
		value := dataAttribute.Value
		if value == nil && dataAttribute.Definition != nil {
			value = dataAttribute.Definition.Value
			if value != nil && value.Value == nil {
				value.updateEnumOrdValue(s._scl.DataTypeTemplates)
			}
		}

		if value != nil {
			s.printValue(daName, dataAttribute, value)
		}
	}
}

func (s *StaticModelGenerator) printValue(daName string, dataAttribute *DataAttribute, value *DataModelValue) {

	// Add a new line
	s.initializerBuffer.WriteString("\n")

	// Conditional initialization
	if s.initializeOnce {
		s.initializerBuffer.WriteString(fmt.Sprintf("if (%s.mmsValue == nil) {\n", daName))
	}

	s.initializerBuffer.WriteString(fmt.Sprintf("%s.mmsValue = ", daName))

	// Handle different data types
	switch dataAttribute.AttributeType {
	case Enumerated, Int8, Int16, Int32, Int64:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newIntegerFromInt32(%d);", cast.ToInt(value.Value)))
	case Int8U, Int16U, Int24U, Int32U:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newUnsignedFromUint32(%d);", cast.ToInt64(value.Value)))
	case Boolean:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newBoolean(%v);", cast.ToBool(value.Value)))
	case OctetString64:
		daValName := daName + "__val"
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newOctetString(0, 64);\n"))
		s.initializerBuffer.WriteString(fmt.Sprintf("uint8_t %s[] = ", daValName))
		appendHexArrayString(s.initializerBuffer, value.Value.([]byte))
		s.initializerBuffer.WriteString(fmt.Sprintf(";\nMmsValue_setOctetString(%s.mmsValue, %s, %d);\n", daName, daValName, len(value.Value.([]byte))))
	case CodedEnum:
		s.initializerBuffer.WriteString("MmsValue_newBitString(2);\n")
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_setBitStringFromIntegerBigEndian(%s.mmsValue, %v);\n", daName, value.Value))
	case UnicodeString255:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newMmsString(\"%s\");", value.Value))
	case VisibleString32, VisibleString64, VisibleString129, VisibleString255, VisibleString65, Currency:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newVisibleString(\"%s\");", value.Value))
	case Float32:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newFloat(%v);", value.Value))
	case Float64:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newDouble(%v);", value.Value))
	case Timestamp:
		s.initializerBuffer.WriteString(fmt.Sprintf("MmsValue_newUtcTimeByMsTime(%v);", value.Value))
	default:
		fmt.Printf("Unknown default value for %s type: %s\n", daName, dataAttribute.AttributeType.ToString())
		s.initializerBuffer.WriteString("NULL;")
	}

	s.initializerBuffer.WriteString("\n")

	// Close initialization block if necessary
	if s.initializeOnce {
		s.initializerBuffer.WriteString("}\n")
	}
}

// appendHexArrayString appends a byte array as a hex string to the builder
func appendHexArrayString(builder *strings.Builder, byteArray []byte) {
	builder.WriteString("{")
	for i, b := range byteArray {
		if i == 0 {
			builder.WriteString(fmt.Sprintf("0x%02X", b))
		} else {
			builder.WriteString(fmt.Sprintf(", 0x%02X", b))
		}
	}
	builder.WriteString("}")
}

func (s *StaticModelGenerator) printDataSets() error {
	logicalDevices := s.accessPoint.Server.LogicalDevices

	for _, logicalDevice := range logicalDevices {
		for _, logicalNode := range logicalDevice.LogicalNodes {
			for _, dataSet := range logicalNode.DataSets {
				dataSetVariableName := fmt.Sprintf("%sds_%s_%s_%s", s.modelPrefix, logicalDevice.Inst, logicalNode.GetName(), dataSet.Name)
				s.dataSetNames = append(s.dataSetNames, dataSetVariableName)
			}
		}
	}

	s.cOut.printlnNone()
	for _, dataSetName := range s.dataSetNames {
		s.cOut.println("extern DataSet %s;", dataSetName)
	}
	s.cOut.printlnNone()

	dataSetNameListIndex := 0
	for _, logicalDevice := range logicalDevices {
		for _, logicalNode := range logicalDevice.LogicalNodes {
			for _, dataSet := range logicalNode.DataSets {

				dataSetVariableName := s.dataSetNames[dataSetNameListIndex]
				dataSetNameListIndex++
				//fcdaCount := 0

				numberOfFcdas := len(dataSet.FCDA)

				s.cOut.printlnNone()

				for i := range dataSet.FCDA {
					dataSetEntryName := dataSetVariableName + "_fcda" + strconv.Itoa(i)
					s.cOut.println("extern DataSetEntry %s;", dataSetEntryName)
					//fcdaCount++
				}
				s.cOut.printlnNone()
				//fcdaCount = 0

				for i, fcda := range dataSet.FCDA {
					dataSetEntryName := dataSetVariableName + "_fcda" + strconv.Itoa(i)
					s.cOut.println("DataSetEntry %s = {", dataSetEntryName)
					s.cOut.println("  \"%s\",", fcda.LdInst)

					mmsVariableNameBuilder := strings.Builder{}
					if fcda.Prefix != "" {
						mmsVariableNameBuilder.WriteString(fcda.Prefix)
					}

					mmsVariableNameBuilder.WriteString(fcda.LnClass)
					if fcda.LdInst != "" {
						mmsVariableNameBuilder.WriteString(fcda.LdInst)
					}
					mmsVariableNameBuilder.WriteString("$" + fcda.Fc)
					mmsVariableNameBuilder.WriteString("$" + toMmsString(fcda.DoName))

					if fcda.DaName != "" {
						mmsVariableNameBuilder.WriteString("$" + toMmsString(fcda.DaName))
					}
					mmsVariableName := mmsVariableNameBuilder.String()

					// Handle array index and component
					variableName := mmsVariableName
					arrayIndex := -1
					var err error
					componentName := ""
					if arrayStart := strings.Index(mmsVariableName, "("); arrayStart != -1 {
						arrayEnd := strings.Index(mmsVariableName, ")")
						arrayIndexStr := mmsVariableName[arrayStart+1 : arrayEnd]

						arrayIndex, err = cast.ToIntE(arrayIndexStr)
						if err != nil {
							return err
						}

						componentNamePart := mmsVariableName[arrayEnd+1:]
						if len(componentNamePart) > 0 && componentNamePart[0] == '$' {
							componentNamePart = componentNamePart[1:]

							if componentNamePart != "" {
								componentName = componentNamePart
							}
						}
					}

					s.cOut.println("  false,")
					s.cOut.println("  \"%s\", ", variableName)
					s.cOut.println("  %d,", arrayIndex)
					if componentName == "" {
						s.cOut.println("  NULL,")
					} else {
						s.cOut.println("  \"%s\",", componentName)
					}
					s.cOut.println("  NULL,")

					if i+1 < numberOfFcdas {
						s.cOut.println("  &%s_fcda%d", dataSetVariableName, i+1)
					} else {
						s.cOut.println("  NULL")
					}
					s.cOut.println("};\n")

				}

				s.cOut.println("DataSet %s = {", dataSetVariableName)
				s.cOut.println("  \"%s\",", logicalDevice.Inst)
				s.cOut.println("  \"%s$%s\",", logicalNode.GetName(), dataSet.Name)
				s.cOut.println("  %d,", numberOfFcdas)
				s.cOut.println("  &%s_fcda0,", dataSetVariableName)

				if dataSetNameListIndex < len(s.dataSetNames) {
					s.cOut.println("  &%s", s.dataSetNames[dataSetNameListIndex])
				} else {
					s.cOut.println("  NULL")
				}
				s.cOut.println("};")
			}
		}
	}

	return nil
}

// createLNSubVariableList contains createReportVariableList,createLogControlVariableList
func (s *StaticModelGenerator) createLNSubVariableList(logicalDevices []*LogicalDevice) {
	for _, ld := range logicalDevices {
		for _, ln := range ld.LogicalNodes {

			rcbCount := 0
			for _, rcb := range ln.ReportControlBlocks {
				maxInstances := 1

				if rcb.RptEnabled != nil {
					maxInstances = rcb.RptEnabled.Max
				}

				for i := 0; i < maxInstances; i++ {
					rcbVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_report" + strconv.Itoa(rcbCount)
					s.rcbVariableNames = append(s.rcbVariableNames, rcbVariableName)
					rcbCount++
				}
			}

			for i := range ln.LogControlBlocks {
				lcbVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_lcb" + strconv.Itoa(i)
				s.lcbVariableNames = append(s.lcbVariableNames, lcbVariableName)
			}

			for i := range ln.Logs {
				logVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_log" + strconv.Itoa(i)
				s.logVariableNames = append(s.logVariableNames, logVariableName)
			}

			for i := range ln.GSEControlBlocks {
				gseVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_gse" + strconv.Itoa(i)
				s.gseVariableNames = append(s.gseVariableNames, gseVariableName)
			}

			for i := range ln.SMVControlBlocks {
				smvVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_smv" + strconv.Itoa(i)
				s.smvVariableNames = append(s.smvVariableNames, smvVariableName)
			}

			for range ln.SettingGroupControlBlocks {
				sgcbVariableName := s.modelPrefix + "_" + ld.Inst + "_" + ln.GetName() + "_sgcb"
				s.sgcbVariableNames = append(s.sgcbVariableNames, sgcbVariableName)
			}
		}
	}

}

func toMmsString(iecString string) string {
	return strings.Replace(iecString, ".", "$", -1)
}
