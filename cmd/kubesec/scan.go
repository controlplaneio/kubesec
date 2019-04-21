package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/sublimino/kubesec/pkg/ruler"
	"github.com/sublimino/kubesec/pkg/server"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type ScanFailedValidationError struct {
}

func (e *ScanFailedValidationError) Error() string {
	return fmt.Sprintf("Kubesec scan failed")
}

var debug bool

func init() {
	scanCmd.Flags().BoolVar(&debug, "debug", false, "turn on debug logs")
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

		if debug {
			z, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}
			logger = z.Sugar()
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

		reports := make([]ruler.Report, 0)
		isJson := json.Valid(fileBytes)
		if isJson {
			report := ruler.NewRuleset(logger).Run(fileBytes)
			reports = append(reports, report)
		} else {
			bits := bytes.Split(fileBytes, []byte(detectLineBreak(fileBytes)+"---"+detectLineBreak(fileBytes)))
			for _, doc := range bits {
				if len(doc) > 0 {
					data, err := yaml.YAMLToJSON(doc)
					if err != nil {
						return err
					}

					report := ruler.NewRuleset(logger).Run(data)
					reports = append(reports, report)

				}
			}
		}

		var lowScore bool
		for _, r := range reports {
			if r.Score <= 0 {
				lowScore = true
				break
			}
		}

		if len(reports) > 1 {
			res, err := json.Marshal(reports)
			if err != nil {
				return err
			}
			fmt.Println(server.PrettyJSON(res))
		} else {
			res, err := json.Marshal(reports[0])
			if err != nil {
				return err
			}
			fmt.Println(server.PrettyJSON(res))
		}

		if len(reports) > 0 && !lowScore {
			return nil
		}

		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true
		return &ScanFailedValidationError{}
	},
}

func detectLineBreak(haystack []byte) string {
  windowsLineEnding := bytes.Contains(haystack, []byte("\r\n"))
  if windowsLineEnding && runtime.GOOS == "windows" {
    return "\r\n"
  }
  return "\n"
}
