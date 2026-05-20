package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_LimitsMemory_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      limits:
        memory: 200Mi
      requests:
        memory: 100Mi
  - name: c2
    resources:
      limits:
        cpu: 100m
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsMemory(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_LimitsMemory_Two_Pods(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      limits:
        memory: 200Mi
      requests:
        memory: 100Mi
  - name: c2
    resources:
      limits:
        memory: 100Mi
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := LimitsMemory(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}
