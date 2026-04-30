package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_RunAsNonRoot_Pod(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
      initContainers:
        - name: init1
        - name: init2
          securityContext:
            runAsNonRoot: false
        - name: init3
          securityContext:
            runAsNonRoot: true
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
	if containers != 3 {
		t.Errorf("Got %v containers wanted %v", containers, 3)
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
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}
