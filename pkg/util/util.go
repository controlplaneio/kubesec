package util

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// NewTabWriter returns a default writer to write tables.
func NewTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.Debug)
}

// PrintTable defines how to print a table.
type PrintTable func(w io.Writer) error

// Print prints any object in yaml, json or table format.
func Print(format string, in interface{}, w io.Writer, fn PrintTable) error {
	var (
		err error
		out []byte
	)

	switch format {
	case "yaml":
		out, err = yaml.Marshal(in)
		if err != nil {
			return err
		}
	case "json":
		out, err = json.MarshalIndent(in, "", "  ")
		if err != nil {
			return err
		}
	case "table":
		if fn == nil {
			return fmt.Errorf("Print table function can not be nil")
		}
		return fn(w)
	default:
		return fmt.Errorf("Unkown printing format: %s", format)
	}

	_, err = fmt.Fprint(w, string(out))
	return err
}
