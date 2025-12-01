package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"unsafe"
)

// DirectoryEntry represents a single directory entry with optional FC annotation
type DirectoryEntry struct {
	Name string
	FC   *FC // Functional Constraint (nil if not applicable)
}

// VariableSpec represents MMS variable specification with type information
type VariableSpec struct {
	Name    string
	Type    MmsType
	IsArray bool
	// For arrays
	ArraySize int
	// For structures
	StructureElements []VariableSpec
}

// ACSIClass represents the ACSI class type for filtering logical node elements
type ACSIClass int

const (
	ACSIClassDataObject ACSIClass = 0
	ACSIClassDataSet    ACSIClass = 1
	ACSIClassBRCB       ACSIClass = 2 // Buffered Report Control Block
	ACSIClassURCB       ACSIClass = 3 // Unbuffered Report Control Block
	ACSIClassLCB        ACSIClass = 4 // Log Control Block
	ACSIClassLog        ACSIClass = 5
	ACSIClassSGCB       ACSIClass = 6  // Setting Group Control Block
	ACSIClassGoCB       ACSIClass = 7  // GOOSE Control Block
	ACSIClassGsCB       ACSIClass = 8  // GSSE Control Block
	ACSIClassMSVCB      ACSIClass = 9  // Multicast Sampled Value Control Block
	ACSIClassUSVCB      ACSIClass = 10 // Unicast Sampled Value Control Block
)

// DataSetInfo represents dataset directory information with deletable flag
type DataSetInfo struct {
	Members     []string // Dataset member references (format: LDName/LNodeName.item[FC])
	IsDeletable bool     // Whether the dataset can be deleted by clients
}

// GetServerDirectory returns list of Logical Devices or files from the server.
// This is the modern, recommended API that replaces GetLogicalDeviceList().
//
// When getFileNames is false: Returns list of Logical Device names (clean, no verbose output)
// When getFileNames is true: Returns list of available files on the server
//
// This function uses the GetServerDirectory ACSI service and does not trigger
// verbose "Data set: xxx (not deletable)" output like the deprecated
// IedConnection_getLogicalDeviceList does.
func (c *Client) GetServerDirectory(getFileNames bool) ([]string, error) {
	var clientError C.IedClientError

	linkedList := C.IedConnection_getServerDirectory(c.conn, &clientError, C.bool(getFileNames))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetLogicalDeviceNames is a convenience wrapper around GetServerDirectory
// that returns logical device names. This is the recommended replacement
// for the deprecated GetLogicalDeviceList() function.
func (c *Client) GetLogicalDeviceNames() ([]string, error) {
	return c.GetServerDirectory(false)
}

// GetLogicalDeviceDirectory returns the directory of logical nodes in a logical device.
// This directly calls IedConnection_getLogicalDeviceDirectory() to get the list of
// logical nodes without triggering verbose "Data set: xxx (not deletable)" output.
//
// Example:
//
//	lnNames, err := client.GetLogicalDeviceDirectory("myLD")
//	// Returns: ["LLN0", "GGIO1", "MMXU1", ...]
func (c *Client) GetLogicalDeviceDirectory(ldName string) ([]string, error) {
	var clientError C.IedClientError

	cLdName := C.CString(ldName)
	defer C.free(unsafe.Pointer(cLdName))

	linkedList := C.IedConnection_getLogicalDeviceDirectory(c.conn, &clientError, cLdName)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetLogicalNodeDirectory returns elements of a logical node filtered by ACSI class.
// This allows retrieving specific types of elements (data objects, datasets, control blocks).
//
// ACSI Classes:
//   - ACSIClassDataObject: Data objects (DO)
//   - ACSIClassDataSet: Data sets
//   - ACSIClassBRCB: Buffered report control blocks
//   - ACSIClassURCB: Unbuffered report control blocks
//   - ACSIClassGoCB: GOOSE control blocks
//   - ACSIClassSGCB: Setting group control blocks
//   - etc.
//
// Example:
//
//	dataObjects, err := client.GetLogicalNodeDirectory("myLD/GGIO1", ACSIClassDataObject)
//	// Returns: ["Ind", "SPCSO1", "Mod", "Beh", "Health"]
//
//	dataSets, err := client.GetLogicalNodeDirectory("myLD/LLN0", ACSIClassDataSet)
//	// Returns: ["Events", "Measurements"]
func (c *Client) GetLogicalNodeDirectory(lnRef string, acsiClass ACSIClass) ([]string, error) {
	var clientError C.IedClientError

	cLnRef := C.CString(lnRef)
	defer C.free(unsafe.Pointer(cLnRef))

	linkedList := C.IedConnection_getLogicalNodeDirectory(c.conn, &clientError, cLnRef, C.ACSIClass(acsiClass))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetDataSetDirectory returns the members of a dataset with deletable information.
// Returns the list of data object references (FCDs/FCDAs) that make up the dataset.
//
// The returned member references use the format: LDName/LNodeName.item[FC]
//
// Example:
//
//	info, err := client.GetDataSetDirectory("myLD/LLN0.Events")
//	// info.Members: ["myLD/GGIO1.Ind.stVal[ST]", "myLD/GGIO1.Ind.q[ST]", "myLD/GGIO1.Ind.t[ST]"]
//	// info.IsDeletable: false (pre-configured) or true (dynamic)
func (c *Client) GetDataSetDirectory(dataSetRef string) (*DataSetInfo, error) {
	var clientError C.IedClientError
	var isDeletable C.bool

	cDataSetRef := C.CString(dataSetRef)
	defer C.free(unsafe.Pointer(cDataSetRef))

	linkedList := C.IedConnection_getDataSetDirectory(c.conn, &clientError, cDataSetRef, &isDeletable)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var members []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			members = append(members, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return &DataSetInfo{
		Members:     members,
		IsDeletable: bool(isDeletable),
	}, nil
}

// GetFileDirectory returns the list of available files on the server
func (c *Client) GetFileDirectory() ([]string, error) {
	return c.GetServerDirectory(true)
}

// GetDataDirectory returns the directory of the given data object (DO) or data attribute.
// Returns list of all data attributes or sub data objects as simple name strings.
//
// This function implements the GetDataDirectory ACSI service.
//
// Example:
//
//	names, err := client.GetDataDirectory("myLD/myLN.myDO")
//	// Returns: ["mag", "q", "t", "subVal"]
func (c *Client) GetDataDirectory(dataRef string) ([]string, error) {
	var clientError C.IedClientError

	cDataRef := C.CString(dataRef)
	defer C.free(unsafe.Pointer(cDataRef))

	linkedList := C.IedConnection_getDataDirectory(c.conn, &clientError, cDataRef)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetDataDirectoryWithFC returns the directory of the given data object with FC annotations.
// Returns list of data attributes or sub data objects with functional constraint appended
// in square brackets (e.g., "mag[MX]", "stVal[ST]", "q[MX]").
//
// This solves the enumeration problem for determining which FC to use when reading values!
//
// Example:
//
//	entries, err := client.GetDataDirectoryWithFC("myLD/myLN.myDO")
//	// Returns: ["mag[MX]", "q[MX]", "t[MX]", "stVal[ST]"]
func (c *Client) GetDataDirectoryWithFC(dataRef string) ([]DirectoryEntry, error) {
	var clientError C.IedClientError

	cDataRef := C.CString(dataRef)
	defer C.free(unsafe.Pointer(cDataRef))

	linkedList := C.IedConnection_getDataDirectoryFC(c.conn, &clientError, cDataRef)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []DirectoryEntry
	current := linkedList
	for current != nil {
		if current.data != nil {
			name := C.GoString((*C.char)(current.data))
			entry := parseDirectoryEntryWithFC(name)
			result = append(result, entry)
		}
		current = current.next
	}

	return result, nil
}

// GetDataDirectoryByFC returns the directory filtered by specific functional constraint.
// More precise than trying multiple FCs - only returns attributes/objects with the given FC.
//
// Example:
//
//	names, err := client.GetDataDirectoryByFC("myLD/myLN.myDO", FC_MX)
//	// Returns only measurement attributes: ["mag", "q", "t"]
func (c *Client) GetDataDirectoryByFC(dataRef string, fc FC) ([]string, error) {
	var clientError C.IedClientError

	cDataRef := C.CString(dataRef)
	defer C.free(unsafe.Pointer(cDataRef))

	linkedList := C.IedConnection_getDataDirectoryByFC(c.conn, &clientError, cDataRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetVariableSpecification returns detailed MMS variable specification for a data attribute.
// This allows you to determine the exact type (float32/float64/int32, etc.) BEFORE reading,
// eliminating guesswork in generic client implementations.
//
// Example:
//
//	spec, err := client.GetVariableSpecification("myLD/myLN.myDO.mag.f", FC_MX)
//	if spec.Type == Float {
//	    // Now you know it's a float and can read it appropriately
//	}
func (c *Client) GetVariableSpecification(dataAttrRef string, fc FC) (*VariableSpec, error) {
	var clientError C.IedClientError

	cDataAttrRef := C.CString(dataAttrRef)
	defer C.free(unsafe.Pointer(cDataAttrRef))

	mmsSpec := C.IedConnection_getVariableSpecification(c.conn, &clientError, cDataAttrRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.MmsVariableSpecification_destroy(mmsSpec)

	spec := parseVariableSpecification(dataAttrRef, mmsSpec)
	return spec, nil
}

// GetLogicalNodeVariables returns all MMS variables that are children of the given logical node.
// Returns raw MMS notation (e.g., "GGIO1$ST$Ind1$stVal") including control block variables.
//
// This is useful for complete enumeration but returns MMS paths, not IEC 61850 references.
func (c *Client) GetLogicalNodeVariables(lnRef string) ([]string, error) {
	var clientError C.IedClientError

	cLnRef := C.CString(lnRef)
	defer C.free(unsafe.Pointer(cLnRef))

	linkedList := C.IedConnection_getLogicalNodeVariables(c.conn, &clientError, cLnRef)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetLogicalDeviceVariables returns all MMS variables in a logical device.
// Returns MMS notation (e.g., "GGIO1$ST$Ind1$stVal") and includes control blocks.
//
// Useful for getting a complete picture of all variables in an LD.
func (c *Client) GetLogicalDeviceVariables(ldName string) ([]string, error) {
	var clientError C.IedClientError

	cLdName := C.CString(ldName)
	defer C.free(unsafe.Pointer(cLdName))

	linkedList := C.IedConnection_getLogicalDeviceVariables(c.conn, &clientError, cLdName)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// GetLogicalDeviceDataSets returns all data set names in a logical device.
// Returns MMS notation (e.g., "LLN0$dataset1").
func (c *Client) GetLogicalDeviceDataSets(ldName string) ([]string, error) {
	var clientError C.IedClientError

	cLdName := C.CString(ldName)
	defer C.free(unsafe.Pointer(cLdName))

	linkedList := C.IedConnection_getLogicalDeviceDataSets(c.conn, &clientError, cLdName)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.LinkedList_destroy(linkedList)

	var result []string
	current := linkedList
	for current != nil {
		if current.data != nil {
			result = append(result, C.GoString((*C.char)(current.data)))
		}
		current = current.next
	}

	return result, nil
}

// Helper function to parse FC from name like "mag[MX]"
func parseDirectoryEntryWithFC(name string) DirectoryEntry {
	entry := DirectoryEntry{Name: name}

	// Look for FC in square brackets
	if len(name) > 3 && name[len(name)-1] == ']' {
		for i := len(name) - 2; i >= 0; i-- {
			if name[i] == '[' {
				fcStr := name[i+1 : len(name)-1]
				entry.Name = name[:i]
				fc := parseFCString(fcStr)
				entry.FC = &fc
				break
			}
		}
	}

	return entry
}

// Helper function to parse FC string to FC enum
func parseFCString(fcStr string) FC {
	switch fcStr {
	case "ST":
		return ST
	case "MX":
		return MX
	case "SP":
		return SP
	case "SV":
		return SV
	case "CF":
		return CF
	case "DC":
		return DC
	case "SG":
		return SG
	case "SE":
		return SE
	case "SR":
		return SR
	case "OR":
		return OR
	case "BL":
		return BL
	case "EX":
		return EX
	case "CO":
		return CO
	case "US":
		return US
	case "MS":
		return MS
	case "RP":
		return RP
	case "BR":
		return BR
	case "LG":
		return LG
	default:
		return ST // Default fallback
	}
}

// Helper function to parse MmsVariableSpecification into VariableSpec
func parseVariableSpecification(name string, mmsSpec *C.MmsVariableSpecification) *VariableSpec {
	if mmsSpec == nil {
		return nil
	}

	spec := &VariableSpec{
		Name: name,
		Type: MmsType(C.MmsVariableSpecification_getType(mmsSpec)),
	}

	// Handle integer subtypes (Int8, Int16, Int32, Int64)
	if spec.Type == Integer {
		size := int(mmsSpec.typeSpec[0])
		switch size {
		case 8:
			spec.Type = Int8
		case 16:
			spec.Type = Int16
		case 32:
			spec.Type = Int32
		default:
			spec.Type = Int64
		}
	}

	// Handle unsigned subtypes
	if spec.Type == Unsigned {
		size := int(mmsSpec.typeSpec[0])
		switch size {
		case 8:
			spec.Type = Uint8
		case 16:
			spec.Type = Uint16
		default:
			spec.Type = Uint32
		}
	}

	// Handle arrays
	if spec.Type == Array {
		spec.IsArray = true
		spec.ArraySize = int(C.MmsVariableSpecification_getSize(mmsSpec))
	}

	// Handle structures (could be expanded to parse child elements)
	if spec.Type == Structure {
		// For now, just mark as structure
		// Full structure parsing would require recursive calls
		spec.StructureElements = []VariableSpec{}
	}

	return spec
}
