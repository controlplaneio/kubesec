package main

import (
	"encoding/json"
	"fmt"
	"github.com/controlplaneio/kubesec/pkg/ruler"
	"github.com/controlplaneio/kubesec/pkg/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
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

		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		reports, err := ruler.NewRuleset(logger).Run(fileBytes)
		if err != nil {
			return err
		}

		if len(reports) == 0 {
			return fmt.Errorf("invalid input %s", filename)
		}

		var lowScore bool
		for _, r := range reports {
			if r.Score <= 0 {
				lowScore = true
				break
			}
		}

		// Sort reports into custom order
		for _, r := range reports {
			if r.Valid {
				if len(r.Scoring.Critical) > 1 {
					sort.Sort(ruler.RuleRefCustomOrder(r.Scoring.Critical))
				}
				if len(r.Scoring.Advise) > 1 {
					sort.Sort(ruler.RuleRefCustomOrder(r.Scoring.Advise))
				}
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
