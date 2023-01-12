package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/controlplaneio/kubesec/v2/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	var format string
	var printRulesCmd = &cobra.Command{
		Use:   `print-rules`,
		Short: "Print all the scanning rules with their associated scores",
		Example: `  kubesec print-rules
  kubesec print-rules -f yaml
  kubesec print-rules -f table`,
	}

	printRulesCmd.Flags().StringVarP(&format, "format", "f", "json", "Set output format (json, yaml, table)")
	printRulesCmd.RunE = func(cmd *cobra.Command, args []string) error {
		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		ruleSet := ruler.NewRuleset(logger)

		// Sort by rule ID
		sort.Slice(ruleSet.Rules, func(i, j int) bool {
			return ruleSet.Rules[i].ID < ruleSet.Rules[j].ID
		})

		printTableFn := func(w io.Writer) error {
			tw := util.NewTabWriter(w)
			fmt.Fprintf(tw, "ID\tReason\tPoints\tKinds\n")
			for _, rule := range ruleSet.Rules {
				fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t\n",
					rule.ID,
					rule.Reason,
					rule.Points,
					strings.Join(rule.Kinds, ","),
				)
			}
			return tw.Flush()
		}

		return util.Print(format, ruleSet.Rules, os.Stdout, printTableFn)
	}

	rootCmd.AddCommand(printRulesCmd)
}
