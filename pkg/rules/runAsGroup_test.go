package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_RunAsGroup_Pod(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        runAsGroup: 10001
      initContainers:
        - name: init1
        - name: init2
          securityContext:
            runAsGroup: 0
        - name: init2
          securityContext:
            runAsGroup: 99999
      containers:
        - name: c1
        - name: c2
          securityContext:
            runAsGroup: 999
        - name: c2
          securityContext:
            runAsGroup: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 4 {
		t.Errorf("Got %v containers wanted %v", containers, 4)
	}
}

func Test_RunAsGroup_InitContainers(t *testing.T) {
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
            runAsGroup: 1
        - name: init2
          securityContext:
            runAsGroup: 10001
      containers:
        - name: c1
        - name: c2
          securityContext:
            runAsGroup: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_RunAsGroup_Guid_Less_Than_10000(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsGroup: 999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_RunAsGroup_Pod_Group_99999(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsGroup: 99999
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}
func Test_RunAsGroup_Two_Containers(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsGroup: 99999
  - name: c2
    securityContext:
      runAsGroup: 0
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_RunAsGroup_Three_Containers(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsGroup: 99999
  - name: c2
    securityContext:
      runAsGroup: 12
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_RunAsGroup_Nil(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    securityContext:
      runAsGroup:
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RunAsGroup(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
