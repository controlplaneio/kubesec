package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_CapDropAll_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      capabilities:
        drop:
          - ALL
  - name: c2
    securityContext:
      capabilities:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapDropAll(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_CapDropAllMissing_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      capabilities:
        drop:
  - name: c2
    securityContext:
      capabilities:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapDropAll(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_CapDropAll_InitContainers(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: init1
    securityContext:
      capabilities:
        drop:
          - ALL
  containers:
  - name: c1
    securityContext:
      capabilities:
        drop:
          - ALL
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapDropAll(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_CapDropAll_Missing(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: init1
  containers:
  - name: c1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapDropAll(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
