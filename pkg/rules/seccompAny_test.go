package rules

import (
	"testing"

	"github.com/ghodss/yaml"
)

func Test_SeccompAny(t *testing.T) {
	for _, tc := range testCasesSeccomp {
		t.Run(tc.description, func(t *testing.T) {
			json, err := yaml.YAMLToJSON([]byte(tc.manifest))
			if err != nil {
				t.Fatal(err.Error())
			}

			count := SeccompAny(json)
			expectedCount := 0
			if tc.expectedProfileType == tcprofSeccompRuntimeDefault || tc.expectedProfileType == tcprofSeccompLocalhost {
				expectedCount = 1
			}

			if count != expectedCount {
				t.Errorf("Expected count was %v but received %v", expectedCount, count)
			}
		})
	}
}
