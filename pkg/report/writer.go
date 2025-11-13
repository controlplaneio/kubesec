package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/controlplaneio/kubesec/v2/pkg/ruler"
)

// Heavily based on aquasecurity/trivy's reporter

// Now returns the current time
var Now = time.Now

type reports ruler.Reports

// WriteReports writes the result to output, format as passed in argument
func WriteReports(format string, output io.Writer, reports reports, outputTemplate string) error {
	var writer Writer
	switch format {
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
		return errors.New("unrecognized format specified")
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
