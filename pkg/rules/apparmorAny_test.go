package rules

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/thedevsaddam/gojsonq/v2"
)

func Test_ApparmorAny(t *testing.T) {
	for _, tc := range testCasesApparmor {
		t.Run(tc.description, func(t *testing.T) {
			json, err := yaml.YAMLToJSON([]byte(tc.manifest))
			if err != nil {
				t.Fatal(err.Error())
			}

			count := ApparmorAny(json)
			expectedCount := 0
			if tc.expectedProfileType == tcprofAppArmorRuntimeDefault || tc.expectedProfileType == tcprofAppArmorLocalhost {
				expectedCount = 1
			}

			if count != expectedCount {
				t.Errorf("Expected count was %v but received %v", expectedCount, count)
			}
		})
	}
}

func Test_isApparmorUnconfined(t *testing.T) {
	testCases := []struct {
		description        string
		json               string
		expectedUnconfined bool
		expectedResult     *checkSecurityContextResult
	}{
		{
			description:        "field missing",
			json:               `{}`,
			expectedUnconfined: false,
			expectedResult: &checkSecurityContextResult{
				unset: true,
				valid: false,
			},
		},
		{
			description:        "non-string field",
			json:               `{"securityContext":{"appArmorProfile":{"type":123}}}`,
			expectedUnconfined: false,
			expectedResult: &checkSecurityContextResult{
				unset: true,
				valid: false,
			},
		},
		{
			description:        "Unconfined when expectedUnconfined=true",
			json:               `{"securityContext":{"appArmorProfile":{"type":"Unconfined"}}}`,
			expectedUnconfined: true,
			expectedResult: &checkSecurityContextResult{
				unset: false,
				valid: true,
			},
		},
		{
			description:        "Unconfined when expectedUnconfined=false",
			json:               `{"securityContext":{"appArmorProfile":{"type":"Unconfined"}}}`,
			expectedUnconfined: false,
			expectedResult: &checkSecurityContextResult{
				unset: false,
				valid: false,
			},
		},
		{
			description:        "Profile=RuntimeDefault when expectedUnconfined=false",
			json:               `{"securityContext":{"appArmorProfile":{"type":"RuntimeDefault"}}}`,
			expectedUnconfined: false,
			expectedResult: &checkSecurityContextResult{
				unset: false,
				valid: true,
			},
		},
		{
			description:        "Profile=RuntimeDefault when expectedUnconfined=true",
			json:               `{"securityContext":{"appArmorProfile":{"type":"RuntimeDefault"}}}`,
			expectedUnconfined: true,
			expectedResult: &checkSecurityContextResult{
				unset: false,
				valid: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			jq := gojsonq.New().FromString(tc.json)
			result := isApparmorUnconfined(jq, tc.expectedUnconfined)

			if tc.expectedResult.unset != result.unset {
				t.Errorf("expected 'checkSecurityContextResult.unset' value for test was %v but received %v instead",
					tc.expectedResult.unset, result.unset)
			}

			if tc.expectedResult.valid != result.valid {
				t.Errorf("expected 'checkSecurityContextResult.valid' value for test was %v but received %v instead",
					tc.expectedResult.valid, result.valid)
			}
		})
	}
}
