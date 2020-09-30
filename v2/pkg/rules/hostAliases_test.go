package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_HostAliases(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  hostAliases:
  - ip: "127.0.0.1"
    hostnames:
    - "foo.local"
    - "bar.local"
  - ip: "10.1.2.3"
    hostnames:
    - "foo.remote"
    - "bar.remote"
  containers:
  - name: cat-hosts
    image: busybox
    command:
    - cat
    args:
    - "/etc/hosts"
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostAliases(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_HostAliases_empty(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  hostAliases:
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostAliases(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_HostAliases_Missing(t *testing.T) {
	// from https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := HostAliases(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
