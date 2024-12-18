package cmds

import "testing"

func TestGenmodelCommand(t *testing.T) {
	args := []string{
		"genmodel",
		"complexModel.cid", // Mock ICD file
		"/Users/jefftao/Documents/Code/GitHub/libiec61850-1.5/examples/server_example_write_handler", // Output directory
	}

	command := New()
	command.SetArgs(args)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
}
