package rules

import (
	"testing"

	"github.com/thedevsaddam/gojsonq/v2"
	"sigs.k8s.io/yaml"
)

func testCheckSecurityContextRule(json []byte) int {
	return checkSecurityContext(
		json,
		true,
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			v := jq.From("securityContext.myAttribute").Get()

			res := checkSecurityContextResult{}
			if v == nil {
				res.unset = true
				return res
			}

			if v.(bool) == true {
				res.valid = true
			}

			return res
		})
}

func TestCheckSecurityContext(t *testing.T) {
	var tests = []struct {
		name     string
		data     string
		expected int
	}{
		{
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        myAttribute: true
      initContainers:
        - name: init1
        - name: init2
`,
			expected: 2,
		},
		{
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        myAttribute: false
      initContainers:
        - name: init1
        - name: init2
`,
			expected: 0,
		},
		{
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      initContainers:
        - name: init1
        - name: init2
`,
			expected: 0,
		},
		{
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        myAttribute: true
      initContainers:
        - name: init1
        - name: init2
          securityContext:
            myAttribute: false
        - name: init3
          securityContext:
            myAttribute: true
      containers:
        - name: c1
        - name: c2
          securityContext:
            myAttribute: false
        - name: c3
          securityContext:
            myAttribute: true
        - name: c4
          securityContext:
            myAttribute: true
`,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json, err := yaml.YAMLToJSON([]byte(tt.data))
			if err != nil {
				t.Fatal(err.Error())
			}

			containers := testCheckSecurityContextRule(json)
			if containers != tt.expected {
				t.Errorf("Got %v containers wanted %v", containers, tt.expected)
			}
		})
	}
}
