package rules

import (
	"testing"

	"github.com/ghodss/yaml"
)

func Test_apparmorUnconfined(t *testing.T) {
	for _, tc := range testCasesApparmor {
		t.Run(tc.description, func(t *testing.T) {
			json, err := yaml.YAMLToJSON([]byte(tc.manifest))
			if err != nil {
				t.Fatal(err.Error())
			}

			count := ApparmorUnconfined(json)
			expectedCount := 0
			if tc.expectedProfileType == tcprofAppArmorUnconfined {
				expectedCount = 1
			}

			if count != expectedCount {
				t.Errorf("Expected count was %v but received %v", expectedCount, count)
			}
		})
	}
}
