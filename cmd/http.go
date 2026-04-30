package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/controlplaneio/kubesec/v2/pkg/server"
	"github.com/controlplaneio/kubesec/v2/pkg/util"
	"github.com/spf13/cobra"
)

const (
	envVarPort        = "PORT"
	envVarKubesecAddr = "KUBESEC_ADDR"
)

var keypath string

func init() {
	httpCmd.Flags().StringVarP(&keypath, "keypath", "k", "", "Path to in-toto link signing key")
	httpCmd.Flags().StringVar(&k8sVersion, "kubernetes-version", "", "Kubernetes version to validate manifets")
	httpCmd.Flags().StringSliceVar(&schemaLocations, "schema-location", []string{}, "Override schema location search path, local or http (can be specified multiple times)")

	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   `http [[ip:]port]`,
	Short: "Starts kubesec HTTP server on the specified IP address (optional) and port",
	Example: `  kubesec http 8080
  kubesec http :8080
  kubesec http 127.0.0.1:808
  kubesec http [::1]:8080

  KUBESEC_ADDR=8080 kubesec http
  KUBESEC_ADDR=:8080 kubesec http
  KUBESEC_ADDR=127.0.0.1:8080 kubesec http`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := os.Getenv(envVarKubesecAddr)

		// Keep PORT env var for backward compatibility
		if v := os.Getenv(envVarPort); v != "" {
			addr = v
			logger.Warnf("usage of %s environment variable is depecrated, use %s instead",
				envVarPort, envVarKubesecAddr)
		}

		// CLI args have precedence over environment variables
		if len(args) == 1 {
			addr = args[0]
		}

		if addr == "" {
			return fmt.Errorf("[[ip:]port] is missing, set CLI argument or use %s environment variable",
				envVarKubesecAddr)
		}

		addr, err := util.SanitizeAddr(addr)
		if err != nil {
			return err
		}

		stopCh := server.SetupSignalHandler()
		jsonLogger, err := NewLogger("info", "json")
		if err != nil {
			return fmt.Errorf("unable to create new logger: %w", err)
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

		server.ListenAndServe(addr, time.Minute, jsonLogger, stopCh, keypath, schemaConfig)
		return nil
	},
}
