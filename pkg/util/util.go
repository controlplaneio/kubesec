package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
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
		// yaml.Marshal() was producing yaml with 4 space indentation
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		err := enc.Encode(in)
		if err != nil {
			return err
		}
		err = enc.Close()
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
			return fmt.Errorf("print table function can not be nil")
		}
		return fn(w)
	default:
		return fmt.Errorf("unknown printing format: %s", format)
	}

	_, err = fmt.Fprint(w, string(out))
	return err
}

// Sanitize and validate HTTP listen address.
// It accepts "port", ":port", "ip:port", "[ipv6]:port" as valid addresses
// if only "port" is provided, the semi-colon prefix is added
func SanitizeAddr(addr string) (string, error) {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	if ip != "" {
		if i := net.ParseIP(ip); i == nil {
			return "", fmt.Errorf("invalid IP address: %q", ip)
		}
	}

	if p, err := strconv.ParseInt(port, 10, 32); err != nil || p < 1 || p > 65535 {
		return "", fmt.Errorf("invalid port: %q", port)
	}

	return addr, nil
}
