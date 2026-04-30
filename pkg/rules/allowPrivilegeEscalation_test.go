package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_AllowPrivilegeEscalation(t *testing.T) {
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
            allowPrivilegeEscalation: true
        - name: c2
          securityContext:
            allowPrivilegeEscalation: true
        - name: c3
          securityContext:
            allowPrivilegeEscalation: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := AllowPrivilegeEscalation(json)
	if containers != 3 {
		t.Errorf("Got %v containers wanted %v", containers, 3)
	}
}

func Test_AllowPrivilegeEscalation_InitContainers(t *testing.T) {
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
            allowPrivilegeEscalation: true
        - name: init2
          securityContext:
            allowPrivilegeEscalation: false
        - name: init3
      containers:
        - name: c1
        - name: c2
          securityContext:
            allowPrivilegeEscalation: false
        - name: c3
          securityContext:
            allowPrivilegeEscalation: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := AllowPrivilegeEscalation(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}
