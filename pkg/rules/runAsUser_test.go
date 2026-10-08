package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_RunAsUser_Pod(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        runAsUser: 10001
      initContainers:
        - name: init1
        - name: init2
          securityContext:
            runAsUser: 0
        - name: init2
          securityContext:
            runAsUser: 99999
      containers:
        - name: c1
        - name: c2
          securityContext:
            runAsUser: 999
        - name: c2
          securityContext:
            runAsUser: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsUser(json)
	if containers != 4 {
		t.Errorf("Got %v containers wanted %v", containers, 4)
	}
}

func Test_RunAsUser_InitContainers(t *testing.T) {
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
            runAsUser: 1
        - name: init2
          securityContext:
            runAsUser: 10001
      containers:
        - name: c1
        - name: c2
          securityContext:
            runAsUser: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsUser(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_RunAsUser(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsUser: 999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsUser(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_RunAsUser_User_99999(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsUser: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsUser(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}
