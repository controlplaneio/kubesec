package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints kubesec version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("version %s\ngit commit %s\nbuild date %s\n", version, commit, date)
		return nil
	},
}
