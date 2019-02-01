package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sublimino/kubesec/pkg/server"
	"time"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   `http [port]`,
	Short: "Starts kubesec HTTP server on the specified port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("port is required")
		}

		stopCh := server.SetupSignalHandler()
		server.ListenAndServe(args[0], time.Minute, logger, stopCh)
		return nil
	},
}
