package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/sublimino/kubesec/pkg/ruler"
	"github.com/sublimino/kubesec/pkg/server"
	"io/ioutil"
	"path/filepath"
)

type ScanFailedValidationError struct {
}

func (e *ScanFailedValidationError) Error() string {
	return fmt.Sprintf("Kubesec scan failed")
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:     `scan [file]`,
	Short:   "Scans Kubernetes resource YAML or JSON",
	Example: `  scan ./deployment.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file path is required")
		}

		filename, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		logger.Debugf("scan filename is %v", filename)

		var data []byte
		isJson := json.Valid(fileBytes)
		if isJson {
			data = fileBytes
		} else {
			data, err = yaml.YAMLToJSON(fileBytes)
			if err != nil {
				return err
			}
		}

		report := ruler.NewRuleset(logger).Run(data)
		res, err := json.Marshal(report)
		if err != nil {
			return err
		}

		fmt.Println(server.PrettyJSON(res))
		if report.Score > 0 {
			return nil
		}

		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true
		return &ScanFailedValidationError{}
	},
}
