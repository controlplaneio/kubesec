package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_CapSysAdmin_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      capabilities:
        add:
          - SYS_ADMIN
          - SYS_TIME
  - name: c2
    securityContext:
      capabilities:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapSysAdmin(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_CapSysAdmin_InitContainers(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: init1
    securityContext:
      capabilities:
        add:
          - SYS_ADMIN
  containers:
  - name: c1
    securityContext:
      capabilities:
        add:
          - SYS_ADMIN
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := CapSysAdmin(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_CapSysAdmin_Missing(t *testing.T) {
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

	containers := CapSysAdmin(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
