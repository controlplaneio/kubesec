package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_RunAsNonRoot(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: c1
          securityContext:
            runAsNonRoot: true
        - name: c2
          securityContext:
            runAsNonRoot: true
        - name: c3
          securityContext:
            runAsNonRoot: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsNonRoot(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_RunAsNonRoot_InitContainers(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      initContainers:
        - name: init1
          securityContext:
            runAsNonRoot: true
        - name: init2
          securityContext:
            runAsNonRoot: false
        - name: init3
      containers:
        - name: c1
        - name: c2
          securityContext:
            runAsNonRoot: false
        - name: c3
          securityContext:
            runAsNonRoot: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsNonRoot(json)
	if containers != 4 {
		t.Errorf("Got %v containers wanted %v", containers, 4)
	}
}
