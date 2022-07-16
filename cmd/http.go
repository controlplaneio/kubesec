package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/server"
	"github.com/spf13/cobra"
)

func init() {
	// FIXME: I don't understand why I need a reference to keypath here,
	// and the cobra docs don't make it exactly clear.
	var keypath string
	var schemaDir string
	httpCmd.Flags().StringVarP(&keypath, "keypath", "k", "", "Path to in-toto link signing key")
	httpCmd.Flags().StringVarP(&schemaDir, "schema-dir", "s", "", "Sets the directory for the json schemas")
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
		jsonLogger, _ := NewLogger("info", "json")

		keypath := cmd.Flag("keypath").Value.String()

		schemaDir := cmd.Flag("schema-dir").Value.String()

		server.ListenAndServe(port, time.Minute, jsonLogger, stopCh, keypath, schemaDir)
		return nil
	},
}
