package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wendy512/iec61850/scl"
	"os"
)

var (
	icdFile        string
	ied            string
	ap             string
	outDir         string
	outFileName    string
	modelPrefix    string
	initializeOnce bool
	_scl           *scl.SCL
)

func New() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "scltool",
		Short: "scltool is used to browse models and generate model code and configuration",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			icdFile = args[0]
			var err error

			if err = checkICDFileExists(icdFile, "ICD file does not exist: %s"); err != nil {
				return err
			}

			parser := scl.NewParser(icdFile)
			_scl, err = parser.Parse()
			if err != nil {
				return err
			}

			if len(args) > 1 {
				outDir = args[1]

				if _, err = os.Stat(outDir); err != nil {
					if os.IsNotExist(err) {
						if err = os.MkdirAll(outFileName, os.ModePerm); err != nil {
							return err
						}
					}

					if err != nil {
						return err
					}
				}
				return nil
			}

			return nil
		},
	}

	genmodelCommand := &cobra.Command{
		Use:   "genmodel <ICD file> <Output file directory>",
		Short: "Generate model files from an ICD file",
		Args:  cobra.ExactArgs(2),
		Run:   runGenModel,
	}

	genmodelCommand.Flags().StringVar(&ied, "ied", "", "IED name")
	genmodelCommand.Flags().StringVar(&ap, "ap", "", "AccessPoints name")
	genmodelCommand.Flags().StringVarP(&outFileName, "out", "o", "static_model", "Output name")
	genmodelCommand.Flags().StringVarP(&modelPrefix, "modelPrefix", "m", "iedModel", "Model prefix name")
	genmodelCommand.Flags().BoolVarP(&initializeOnce, "initializeonce", "i", false, "Initialize once")

	rootCommand.AddCommand(genmodelCommand)

	return rootCommand
}

func checkICDFileExists(filePath, template string) error {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(template, filePath)
		}
		return err
	}
	return nil
}

func runGenModel(cmd *cobra.Command, args []string) {
	if err := scl.NewStaticModelGenerator(_scl, ied, ap, outDir, outFileName, modelPrefix, initializeOnce).Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
