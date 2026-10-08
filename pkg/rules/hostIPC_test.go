package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_HostIPC(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostIPC: true
      containers:
        - name: c1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostIPC(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_HostIPC_Disabled(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostIPC: false
      containers:
        - name: c1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostIPC(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_HostIPC_Missing(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: c1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostIPC(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
