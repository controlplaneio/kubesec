package util

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	type TestStruct struct {
		Name   string   `json:"name"  yaml:"name"`
		Values []string `json:"values"  yaml:"values"`
	}

	var TestObject = TestStruct{
		Name:   "test",
		Values: []string{"val1", "val2"},
	}

	var tests = []struct {
		name         string
		format       string
		printTableFn PrintTable
		want         string
		wantErr      bool
	}{
		{
			name:         "Print struct in json",
			format:       "json",
			printTableFn: nil,
			want: `{
  "name": "test",
  "values": [
    "val1",
    "val2"
  ]
}`,
			wantErr: false,
		},
		{
			name:         "Print struct in yaml",
			format:       "yaml",
			printTableFn: nil,
			want: `name: test
values:
- val1
- val2
`,
			wantErr: false,
		},
		{
			name:         "Print struct in table without print function",
			format:       "table",
			printTableFn: nil,
			want:         "",
			wantErr:      true,
		},
		{
			name:   "Print struct in table with print function",
			format: "table",
			printTableFn: func(w io.Writer) error {
				tw := NewTabWriter(w)
				fmt.Fprintf(tw, "ID\tValues\n")
				fmt.Fprintf(tw, "%s\t%s\n",
					TestObject.Name,
					strings.Join(TestObject.Values, ","))
				return tw.Flush()
			},
			want: `ID    |Values
test  |val1,val2
`,
			wantErr: false,
		},
		{
			name:    "Print struct in unsupported format",
			format:  "cue",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer

			err := Print(tt.format, TestObject, &b, tt.printTableFn)
			if err == nil && tt.wantErr {
				t.Errorf("Test should have failed")
			}

			if err != nil && !tt.wantErr {
				t.Error(err)
			}

			if tt.want != b.String() {
				t.Errorf("want:\n%s\n\ngot:\n\n%s", tt.want, b.String())
			}
		})
	}
}

func TestSanitizeAddr(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
		expectErr      bool
	}{
		// valid
		{"8080", ":8080", false},
		{":8080", ":8080", false},
		{"127.0.0.1:8080", "127.0.0.1:8080", false},
		{"[::1]:8080", "[::1]:8080", false},
		// invalid
		{"0", "", true},
		{"100000", "", true},
		{"127.0.0.1", "", true},
		{"::1", "", true},
		{"::1:8080", "", true},
		{"example.com", "", true},
		{"example.com:8080", "", true},
	}

	for _, tt := range tests {
		output, err := SanitizeAddr(tt.input)

		if tt.expectErr && err == nil {
			t.Fatalf("an error should have occured with input %s, got nil", tt.input)
		}

		if tt.expectedOutput != output {
			t.Errorf("unexpected output, want %s, got %s", tt.expectedOutput, output)
		}
	}
}
