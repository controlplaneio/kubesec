package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/controlplaneio/kubesec/v2/pkg/pss"
	"github.com/controlplaneio/kubesec/v2/pkg/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	var (
		debug      bool
		format     string
		profile    string
		k8sVersion string
	)

	var pssScanCmd = &cobra.Command{
		Use:   `pss-scan [file]`,
		Short: "Scan Kubernetes resources (yaml, json) against Pod Security Standards (PSS) profiles",
		Long: `The default scanning configuration is to validate the manifests
against the highly-restrictive "restricted" profile with the latest version.

For more information about Pod Security Standards (PSS) and the profiles,
refer to the official documentation:

* https://kubernetes.io/docs/concepts/security/pod-security-standards`,
		Example: `  kubesec pss-scan
  kubesec pss-scan --profile baseline -f yaml
  kubesec pss-scan --profile baseline --profile-version v1.26`,
	}

	pssScanCmd.Flags().BoolVar(&debug, "debug", false, "Turn on debug logs")
	pssScanCmd.Flags().IntVar(&exitCodeError, "exit-code", 2, "Set the exit-code to use on failure")
	pssScanCmd.Flags().StringVarP(&format, "format", "f", "json", "Set output format (json, yaml)")
	pssScanCmd.Flags().StringVar(&profile, "profile", "restricted", "")
	pssScanCmd.Flags().StringVar(&k8sVersion, "kubernetes-version", "", "Kubernetes version to validate manifets (latest or 1.x)")
	pssScanCmd.RunE = func(cmd *cobra.Command, args []string) error {
		rootCmd.SilenceErrors = true

		if len(args) < 1 {
			return fmt.Errorf("file path is required")
		}

		if format != "json" && format != "yaml" {
			return fmt.Errorf("output format not supported: %s", format)
		}

		if debug {
			z, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}
			logger = z.Sugar()
		}

		rootCmd.SilenceUsage = true

		file, err := util.ReadManifest(args, absolutePath)
		if err != nil {
			return err
		}

		evaluator, err := pss.NewEvaluator(logger)
		if err != nil {
			return err
		}

		switch k8sVersion {
		case "":
			k8sVersion = "latest"
		case "latest":
			// this could be set as default in cmd but it is left empty
			// for consistency with the scan subcommand
		default:
			if !strings.HasPrefix(k8sVersion, "v") {
				k8sVersion = "v" + k8sVersion
			}
		}

		reports, err := evaluator.Run(file.Name, file.Bytes, profile, k8sVersion)
		if err != nil {
			// This check allows setting a different exit error code
			// and printing of the check report.
			e := &pss.ProfileNotSatisfiedError{}
			if !errors.As(err, &e) {
				return err
			}
			exitCode = exitCodeError
			logger.Debug(err)
		}

		return util.Print(format, reports, os.Stdout, func(w io.Writer) error { return nil })
	}

	rootCmd.AddCommand(pssScanCmd)
}
