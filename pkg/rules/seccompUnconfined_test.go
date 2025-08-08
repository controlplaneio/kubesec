package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_SeccompUnconfined(t *testing.T) {
	for _, tc := range testCasesSeccomp {
		t.Run(tc.description, func(t *testing.T) {
			json, err := yaml.YAMLToJSON([]byte(tc.manifest))
			if err != nil {
				t.Fatal(err.Error())
			}

			count := SeccompUnconfined(json)
			expectedCount := 0
			if tc.expectedProfileType == tcprofSeccompUnconfined {
				expectedCount = 1
			}

			if count != expectedCount {
				t.Errorf("Expected count was %v but received %v", expectedCount, count)
			}
		})
	}
}
