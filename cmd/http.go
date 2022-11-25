package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/controlplaneio/kubesec/v2/pkg/server"
	"github.com/spf13/cobra"
)

var keypath string

func init() {
	httpCmd.Flags().StringVarP(&keypath, "keypath", "k", "", "Path to in-toto link signing key")
	httpCmd.Flags().StringVar(&k8sVersion, "kubernetes-version", "", "Kubernetes version to validate manifets")
	httpCmd.Flags().StringSliceVar(&schemaLocations, "schema-location", []string{}, "Override schema location search path, local or http (can be specified multiple times)")

	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   `http [port]`,
	Short: "Starts kubesec HTTP server on the specified port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("port is required")
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = args[0]
		}

		if _, err := strconv.Atoi(port); err != nil {
			port = args[0]
		}

		stopCh := server.SetupSignalHandler()
		jsonLogger, err := NewLogger("info", "json")
		if err != nil {
			return fmt.Errorf("Unable to create new logger: %w", err)
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

		server.ListenAndServe(port, time.Minute, jsonLogger, stopCh, keypath, schemaConfig)
		return nil
	},
}
