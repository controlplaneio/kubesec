package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_LimitsCPU_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      limits:
        cpu: 300m
      requests:
        cpu: 300m
  - name: c2
    resources:
      limits:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsCPU(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_LimitsCPU_Two_Pods(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      limits:
        cpu: 300m
      requests:
        cpu: 300m
  - name: c2
    resources:
      limits:
  - name: c3
    resources:
      limits:
        cpu: 123m
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsCPU(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_LimitsCPU_Pod_Malformed(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      limits:
      requests:
  - name: c2
    resources:
      limits:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsCPU(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_LimitsCPU_Pod_Missing(t *testing.T) {
	var data = `
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    run: daemonset
  name: daemonset
spec:
  replicas: 1
  selector:
    matchLabels:
      run: daemonset
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        run: daemonset
    spec:
      containers:
      - args:
        - arse
        image: arse
        name: daemonset
        resources: {}
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsCPU(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
