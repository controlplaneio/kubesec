package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/controlplaneio/kubesec/v2/pkg/report"
	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ScanFailedValidationError struct {
}

func (e *ScanFailedValidationError) Error() string {
	return "Kubesec scan failed"
}

var (
	debug           bool
	absolutePath    bool
	format          string
	template        string
	k8sVersion      string
	schemaLocations = []string{}
	outputLocation  string
	exitCode        int
	exitCodeError   int
)

func init() {
	scanCmd.Flags().BoolVar(&debug, "debug", false, "Turn on debug logs")
	scanCmd.Flags().BoolVar(&absolutePath, "absolute-path", false, "Use the absolute path for the file name")
	scanCmd.Flags().StringVarP(&format, "format", "f", "json", "Set output format (json, template)")
	scanCmd.Flags().StringVar(&k8sVersion, "kubernetes-version", "", "Kubernetes version to validate manifets")
	scanCmd.Flags().StringSliceVar(&schemaLocations, "schema-location", []string{}, "Override schema location search path, local or http (can be specified multiple times)")
	scanCmd.Flags().StringVarP(&template, "template", "t", "", "Set output template, it will check for a file or read input as the")
	scanCmd.Flags().StringVarP(&outputLocation, "output", "o", "", "Set output location")
	scanCmd.Flags().IntVar(&exitCodeError, "exit-code", 2, "Set the exit-code to use on failure")
	rootCmd.AddCommand(scanCmd)
}

// File holds the name and contents
type File struct {
	fileName  string
	fileBytes []byte
}

func getInput(args []string) (File, error) {
	var file File

	if len(args) == 1 && (args[0] == "-" || args[0] == "/dev/stdin") {
		fileBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return file, err
		}
		file = File{
			fileName:  "STDIN",
			fileBytes: fileBytes,
		}
		return file, nil
	}
	fileName := args[0]
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return file, err
	}
	if absolutePath {
		fileName = filePath
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return file, err
	}
	file = File{
		fileName:  fileName,
		fileBytes: fileBytes,
	}
	return file, nil
}

var scanCmd = &cobra.Command{
	Use:   `scan [file]`,
	Short: "Scan Kubernetes resources (yaml, json) against kubesec rules",
	Example: `  kubesec scan ./deployment.yaml
  cat file.json | kubesec scan -
  helm template -f values.yaml ./chart | kubesec scan /dev/stdin`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file path is required")
		}

		if debug {
			z, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}
			logger = z.Sugar()
		}

		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		file, err := getInput(args)
		if err != nil {
			return err
		}

		ver := os.Getenv("K8S_SCHEMA_VER")
		if ver != "" && k8sVersion == "" {
			k8sVersion = ver
		}

		loc := os.Getenv("SCHEMA_LOCATION")
		if loc != "" && len(schemaLocations) == 0 {
			schemaLocations = strings.Split(loc, ",")
		}

		schemaConfig := ruler.NewDefaultSchemaConfig()
		schemaConfig.Locations = schemaLocations
		schemaConfig.ValidatorOpts.KubernetesVersion = k8sVersion

		reports, err := ruler.NewRuleset(logger).Run(file.fileName, file.fileBytes, schemaConfig)
		if err != nil {
			return err
		}

		if len(reports) == 0 {
			return fmt.Errorf("invalid input %s", file.fileName)
		}

		var lowScore bool
		for _, r := range reports {
			if r.Score <= 0 {
				lowScore = true
				break
			}
		}

		var buff bytes.Buffer
		err = report.WriteReports(format, &buff, reports, template)
		if err != nil {
			return err
		}

		if outputLocation != "" {
			err = ioutil.WriteFile(outputLocation, buff.Bytes(), 0644)
			if err != nil {
				logger.Debugf("Couldn't write output to %s", outputLocation)
			}
		}

		out := buff.String()
		fmt.Println(out)

		if len(reports) > 0 {
			if lowScore {
				exitCode = exitCodeError
			}
			return nil
		}

		return &ScanFailedValidationError{}
	},
}
