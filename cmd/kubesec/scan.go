package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/controlplaneio/kubesec/pkg/ruler"
	"github.com/controlplaneio/kubesec/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

func getInput(args []string) ([]ruler.File, error) {
	var files []ruler.File

	if len(args) == 1 && (args[0] == "-" || args[0] == "/dev/stdin") {
		fileBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return files, err
		}
		file := ruler.File{
			FileName:  "STDIN",
			FileBytes: fileBytes,
		}
		return append(files, file), nil
	}

	for _, arg := range args {
		filename, err := filepath.Abs(arg)
		if err != nil {
			return files, err
		}
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return files, err
		}
		file := ruler.File{
			FileName:  filename,
			FileBytes: fileBytes,
		}
		files = append(files, file)
	}

	return files, nil
}

var scanCmd = &cobra.Command{
	Use:     `scan [files]`,
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

		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		files, err := getInput(args)
		if err != nil {
			return err
		}
		reports, err := ruler.NewRuleset(logger).Run(files)
		if err != nil {
			return err
		}

		if len(reports) == 0 {
			var fileNames []string
			for _, f := range files {
				fileNames = append(fileNames, f.FileName)
			}
			return fmt.Errorf("invalid inputs: %s", strings.Join(fileNames, ", "))
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
