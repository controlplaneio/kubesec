package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
	"github.com/pterm/pterm"
)

// Heavily based on aquasecurity/trivy's reporter

// Now returns the current time
var Now = time.Now

type reports ruler.Reports

// WriteReports writes the result to output, format as passed in argument
func WriteReports(format string, output io.Writer, reports reports, outputTemplate string) error {
	var writer Writer
	switch format {
	case "table":
		writer = &TableWriter{Output: output}
	case "json":
		writer = &JSONWriter{Output: output}
	case "template":
		var err error
		if len(outputTemplate) == 0 {
			return errors.New("template is unset, please specify with --template")
		}
		if writer, err = NewTemplateWriter(output, outputTemplate); err != nil {
			return err
		}
	default:
		return errors.New("Unrecognized format specified")
	}

	if err := writer.Write(reports); err != nil {
		return err
	}
	return nil
}

// Writer defines the result write operation
type Writer interface {
	Write(reports) error
}

// JSONWriter implements result Writer
type JSONWriter struct {
	Output io.Writer
}

// PrettyJSON will indent JSON to be pretty
func PrettyJSON(jsonBytes []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, jsonBytes, "", "  ")
	if err != nil {
		return jsonBytes, err
	}
	return out.Bytes(), nil
}

// Write writes the reports in JSON format
func (jw JSONWriter) Write(reports reports) error {
	output, err := json.Marshal(reports)
	if err != nil {
		return err
	}

	formattedOutput, err := PrettyJSON(output)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprint(jw.Output, string(formattedOutput)); err != nil {
		return err
	}
	return nil
}

// TemplateWriter write result in custom format defined by user's template
type TemplateWriter struct {
	Output   io.Writer
	Template *template.Template
}

// NewTemplateWriter is the factory method to return TemplateWriter object
func NewTemplateWriter(output io.Writer, outputTemplate string) (*TemplateWriter, error) {
	// if outputTemplate is a file read it and use that
	if _, err := os.Stat(outputTemplate); err == nil {
		buf, err := os.ReadFile(outputTemplate)
		if err != nil {
			return nil, err
		}
		outputTemplate = string(buf)
	}

	tmpl, err := template.New("output template").Funcs(template.FuncMap{
		"endWithPeriod": func(input string) string {
			if !strings.HasSuffix(input, ".") {
				input += "."
			}
			return input
		},
		"toLower": func(input string) string {
			return strings.ToLower(input)
		},
		"escapeString": func(input string) string {
			return html.EscapeString(input)
		},
		"getCurrentTime": func() string {
			return Now().UTC().Format(time.RFC3339Nano)
		},
		"joinSlices": func(slices ...[]ruler.RuleRef) []ruler.RuleRef {
			var resultSlice []ruler.RuleRef
			for _, slice := range slices {
				resultSlice = append(resultSlice, slice...)
			}
			return resultSlice
		},
	}).Parse(outputTemplate)
	if err != nil {
		return nil, err
	}
	return &TemplateWriter{Output: output, Template: tmpl}, nil
}

// Write writes result
func (tw TemplateWriter) Write(reports reports) error {
	err := tw.Template.Execute(tw.Output, reports)
	if err != nil {
		return err
	}
	return nil
}

// TableWriter implements result Writer for the table format
type TableWriter struct {
	Output io.Writer
}

// Write writes results as a table
func (tw TableWriter) Write(reports reports) error {
	pterm.SetDefaultOutput(tw.Output)

	// Header and Overall Summary Table
	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithMargin(1).
		Println("🛡️  Kubesec Scan Report")

	var summaryData [][]string
	summaryData = append(summaryData, []string{"File", "Object", "Valid", "Score", "# Critical Rule Hits", "# Advise Rule Hits", "# Passed Rule Hits"})

	for _, r := range reports {
		isValid := pterm.LightGreen("✅")
		if !r.Valid {
			isValid = pterm.LightRed("❌")
		}

		summaryData = append(summaryData, []string{
			r.FileName,
			r.Object,
			isValid,
			pterm.Cyan(strconv.Itoa(r.Score)),
			pterm.LightRed(strconv.Itoa(len(r.Scoring.Critical))),
			pterm.LightYellow(strconv.Itoa(len(r.Scoring.Advise))),
			pterm.LightGreen(strconv.Itoa(len(r.Scoring.Passed))),
		})
	}

	// Render the summary table
	err := pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithHeaderStyle(pterm.NewStyle(pterm.FgCyan, pterm.Bold)).
		WithData(summaryData).
		Render()
	if err != nil {
		return err
	}
	pterm.Println()

	// Detailed Tables for Each Report
	for _, r := range reports {
		// Skip if there are no rules to display for this report
		if len(r.Scoring.Critical) == 0 && len(r.Scoring.Advise) == 0 && len(r.Scoring.Passed) == 0 {
			continue
		}

		// Section header for the specific resource
		pterm.DefaultSection.Printf("Details for %s", pterm.LightCyan(r.Object))
		pterm.Printf("File: %s\n", pterm.Gray(r.FileName))
		if r.Message != "" {
			pterm.Info.Println(r.Message)
		}

		var detailData [][]string
		detailData = append(detailData, []string{"Severity", "Rule ID", "Selector", "Reason", "Points"})

		// Helper function to append rules to the table
		appendRules := func(rules []ruler.RuleRef, severityText string) {
			selector := ""

			for _, rule := range rules {
				// split long selectors into multiple lines for table to render correctly
				selector = strings.ReplaceAll(rule.Selector, " | ", " |\n")

				detailData = append(detailData, []string{
					severityText,
					rule.ID,
					selector,
					rule.Reason,
					fmt.Sprintf("%d", rule.Points),
				})
			}
		}

		// Append all rules grouped by severity
		appendRules(r.Scoring.Critical, pterm.LightRed("🔴 Critical"))
		appendRules(r.Scoring.Advise, pterm.LightYellow("🟡 Advise"))
		appendRules(r.Scoring.Passed, pterm.LightGreen("🟢 Passed"))

		// Render the detailed table
		err = pterm.DefaultTable.
			WithHasHeader().
			WithBoxed(true).
			WithRowSeparator("-").
			WithHeaderStyle(pterm.NewStyle(pterm.FgWhite, pterm.Bold)).
			WithData(detailData).
			Render()
		if err != nil {
			return err
		}
		pterm.Println()
	}

	return nil
}
