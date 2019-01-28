// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
  "fmt"
  "os"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var cfgFile string

var buildStamp = "unknown"
var gitHash = "unknown"
var buildVersion = "unknown"

var (
  isVersionFlag = false
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
  Use:   "kubesec",
  Short: "Validate Kubernetes resource security policies",
  Long: `Security
Security
Security
`,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    kubeSecCheck(cmd, args)
  },

  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}

func kubeSecCheck(cmd *cobra.Command, args []string) {

  if isVersionFlag {
    fmt.Fprintf(os.Stderr, "%s %s (build %s %s)\n", os.Args[0], buildVersion, gitHash, buildStamp)
    os.Exit(0)
  }

  kubeResourceYaml := args[0]

  fmt.Println("Resource " + kubeResourceYaml)

  getJsonFromKubeResource()
  getRulesFromFile()
  runRulesOnKubeResource()
  outputResults()
}

func getJsonFromKubeResource() {
  fmt.Println("getJsonFromKubeResource()")
}

func getRulesFromFile() {
  fmt.Println("getRulesFromFile()")
}

func runRulesOnKubeResource() {
  fmt.Println("runRulesOnKubeResource()")
}

func outputResults() {
  fmt.Println("outputResults()")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".kubesec" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".kubesec")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}
