package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var (
	logger *zap.SugaredLogger

	// vars injected by goreleaser at build time
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "kubesec",
	Short: "kubesec command line",
	Long: `
Validate Kubernetes resource security policies`,
}

func main() {
	var err error

	// logger writes to stderr
	zlog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = zlog.Sugar()
	defer logger.Sync()

	rootCmd.SetArgs(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}
