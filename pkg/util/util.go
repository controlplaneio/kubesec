package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
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
		out, err = yaml.Marshal(in)
		if err != nil {
			return err
		}
	case "json":
		out, err = json.MarshalIndent(in, "", "  ")
		if err != nil {
			return err
		}

		// Add a newline for readibility as yaml does
		newline := "\n"
		out = append(out, newline...)
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

// File holds the name and contents.
type File struct {
	Name  string
	Bytes []byte
}

// ReadManifest reads the content of a manifest from either a file or /dev/stdin.
func ReadManifest(args []string, absolutePath bool) (File, error) {
	var file File

	if len(args) == 1 && (args[0] == "-" || args[0] == "/dev/stdin") {
		fileBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return file, err
		}
		file = File{
			Name:  "STDIN",
			Bytes: fileBytes,
		}
		return file, nil
	}
	fileName := args[0]
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return file, err
	}
	if absolutePath {
		fileName = filePath
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return file, err
	}

	// Ensure we're dealing with text files
	contentType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(contentType, "text/plain") {
		return file, fmt.Errorf("Provided file Content-Type is not of supported type text/plain, got: %s", contentType)
	}

	// Remove empty content
	if len(fileBytes) == 0 {
		return file, fmt.Errorf("Provided file is empty")
	}
	file = File{
		Name:  fileName,
		Bytes: fileBytes,
	}
	return file, nil
}

// GetObjectsFromManifest returns a slice of YAML or JSON objects from a manifest content.
func GetObjectsFromManifest(f []byte) [][]byte {
	objects := make([][]byte, 0)
	switch {
	// Find if it is a JSON array (starts with [)
	case f[0] == 91:
		jsonObjects := make([]json.RawMessage, 0)
		err := json.Unmarshal(f, &jsonObjects)
		if err != nil {
			// not JSON or misformated
			return nil
		}

		for _, obj := range jsonObjects {
			objects = append(objects, obj)
		}
	// Consider the rest as YAML or JSON singleton
	default:
		lineBreak := detectLineBreak(f)
		objects = bytes.Split(f, []byte(lineBreak+"---"+lineBreak))
	}

	return objects
}

func detectLineBreak(haystack []byte) string {
	windowsLineEnding := bytes.Contains(haystack, []byte("\r\n"))
	if windowsLineEnding && goruntime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
