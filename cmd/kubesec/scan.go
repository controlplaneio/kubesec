package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/controlplaneio/kubesec/v2/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ScanFailedValidationError struct {
}

func (e *ScanFailedValidationError) Error() string {
	return "Kubesec scan failed"
}

var debug bool

func init() {
	scanCmd.Flags().BoolVar(&debug, "debug", false, "turn on debug logs")
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
	Short: "Scans Kubernetes resource YAML or JSON",
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

		reports, err := ruler.NewRuleset(logger).Run(file.fileName, file.fileBytes)
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

		res, err := json.Marshal(reports)
		if err != nil {
			return err
		}
		fmt.Println(server.PrettyJSON(res))

		if len(reports) > 0 && !lowScore {
			return nil
		}

		return &ScanFailedValidationError{}
	},
}
