package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_ServiceAccountName(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: build-robot
  automountServiceAccountToken: false
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ServiceAccountName(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_ServiceAccountName_empty(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: 
  automountServiceAccountToken: false
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ServiceAccountName(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_ServiceAccountName_Missing(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  automountServiceAccountToken: false
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ServiceAccountName(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
