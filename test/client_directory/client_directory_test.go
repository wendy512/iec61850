package client_directory

import (
	"fmt"
	"testing"

	"github.com/boeboe/iec61850"
	"github.com/boeboe/iec61850/test"
)

// TestGetServerDirectory demonstrates the modern API for getting logical device names
// This replaces the deprecated GetLogicalDeviceList() which causes verbose output
func TestGetServerDirectory(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	// Get logical device names (clean, no verbose output)
	devices, err := client.GetLogicalDeviceNames()
	if err != nil {
		t.Fatalf("get logical device names error: %v\n", err)
	}

	t.Logf("Found %d logical devices:", len(devices))
	for i, device := range devices {
		t.Logf("  [%d] %s", i, device)
	}
}

// TestGetDataDirectoryWithFC demonstrates getting data attributes with FC annotations
// This solves the enumeration problem - you know which FC to use!
func TestGetDataDirectoryWithFC(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	// First get the logical devices
	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	// Example: Get data directory with FC for first device's first logical node
	// Adjust this reference based on your actual server model
	dataRef := fmt.Sprintf("%s/LLN0.Beh", devices[0])

	t.Logf("\nGetting data directory with FC for: %s", dataRef)
	entries, err := client.GetDataDirectoryWithFC(dataRef)
	if err != nil {
		t.Logf("Note: %v (this is expected if the reference doesn't exist)", err)
		return
	}

	t.Logf("Found %d attributes with FC annotations:", len(entries))
	for _, entry := range entries {
		if entry.FC != nil {
			t.Logf("  %s [%v]", entry.Name, *entry.FC)
		} else {
			t.Logf("  %s [no FC]", entry.Name)
		}
	}
}

// TestGetDataDirectoryByFC demonstrates filtering by specific functional constraint
func TestGetDataDirectoryByFC(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	// Example: Get only status (ST) attributes
	dataRef := fmt.Sprintf("%s/LLN0.Beh", devices[0])

	t.Logf("\nGetting ST attributes for: %s", dataRef)
	names, err := client.GetDataDirectoryByFC(dataRef, iec61850.ST)
	if err != nil {
		t.Logf("Note: %v (this is expected if the reference doesn't exist)", err)
		return
	}

	t.Logf("Found %d ST attributes:", len(names))
	for _, name := range names {
		t.Logf("  %s", name)
	}
}

// TestGetVariableSpecification demonstrates type introspection before reading
// This eliminates guesswork - you know the exact type!
func TestGetVariableSpecification(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	// Example: Get variable specification for a data attribute
	// Adjust this reference based on your actual server model
	attrRef := fmt.Sprintf("%s/LLN0.Beh.stVal", devices[0])

	t.Logf("\nGetting variable specification for: %s", attrRef)
	spec, err := client.GetVariableSpecification(attrRef, iec61850.ST)
	if err != nil {
		t.Logf("Note: %v (this is expected if the reference doesn't exist)", err)
		return
	}

	t.Logf("Variable specification:")
	t.Logf("  Name: %s", spec.Name)
	t.Logf("  Type: %v", spec.Type)
	t.Logf("  IsArray: %v", spec.IsArray)
	if spec.IsArray {
		t.Logf("  ArraySize: %d", spec.ArraySize)
	}
}

// TestGetLogicalNodeVariables demonstrates getting all MMS variables in a logical node
func TestGetLogicalNodeVariables(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	// Example: Get all variables in LLN0
	lnRef := fmt.Sprintf("%s/LLN0", devices[0])

	t.Logf("\nGetting all MMS variables for: %s", lnRef)
	variables, err := client.GetLogicalNodeVariables(lnRef)
	if err != nil {
		t.Logf("Note: %v", err)
		return
	}

	t.Logf("Found %d MMS variables (showing first 10):", len(variables))
	for i, variable := range variables {
		if i >= 10 {
			t.Logf("  ... (%d more)", len(variables)-10)
			break
		}
		t.Logf("  %s", variable)
	}
}

// TestGetLogicalDeviceVariables demonstrates getting all MMS variables in a logical device
func TestGetLogicalDeviceVariables(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	t.Logf("\nGetting all MMS variables for LD: %s", devices[0])
	variables, err := client.GetLogicalDeviceVariables(devices[0])
	if err != nil {
		t.Logf("Note: %v", err)
		return
	}

	t.Logf("Found %d MMS variables (showing first 10):", len(variables))
	for i, variable := range variables {
		if i >= 10 {
			t.Logf("  ... (%d more)", len(variables)-10)
			break
		}
		t.Logf("  %s", variable)
	}
}

// TestGetLogicalDeviceDataSets demonstrates getting all dataset names in a logical device
func TestGetLogicalDeviceDataSets(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	devices, err := client.GetLogicalDeviceNames()
	if err != nil || len(devices) == 0 {
		t.Skip("No logical devices available for testing")
	}

	t.Logf("\nGetting all datasets for LD: %s", devices[0])
	datasets, err := client.GetLogicalDeviceDataSets(devices[0])
	if err != nil {
		t.Logf("Note: %v", err)
		return
	}

	t.Logf("Found %d datasets:", len(datasets))
	for _, dataset := range datasets {
		t.Logf("  %s", dataset)
	}
}
