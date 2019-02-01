package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_ReadOnlyRootFilesystem(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: c1
        - name: c2
          securityContext:
            readOnlyRootFilesystem: false
        - name: c3
          securityContext:
            readOnlyRootFilesystem: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ReadOnlyRootFilesystem(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_ReadOnlyRootFilesystem_InitContainers(t *testing.T) {
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
            readOnlyRootFilesystem: true
        - name: init2
          securityContext:
            readOnlyRootFilesystem: false
        - name: init3
      containers:
        - name: c1
        - name: c2
          securityContext:
            readOnlyRootFilesystem: false
        - name: c3
          securityContext:
            readOnlyRootFilesystem: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ReadOnlyRootFilesystem(json)
	if containers != 4 {
		t.Errorf("Got %v containers wanted %v", containers, 4)
	}
}

func Test_ReadOnlyRootFilesystem_NotSpecified(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  containers:
  - name: sec-ctx-demo
    image: gcr.io/google-samples/node-hello:1.0
    securityContext:
      capabilities:
        add:
          - CHOWN
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ReadOnlyRootFilesystem(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_ReadOnlyRootFilesystem_NoContainers(t *testing.T) {
	var data = `
---
apiVersion: extensions/v1beta1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: kubesec
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ReadOnlyRootFilesystem(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
