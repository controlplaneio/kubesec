package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_SeccompAny_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    other: runtime/default
    seccomp.security.alpha.kubernetes.io/pod: runtime/default
    something: runtime/default
spec:
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SeccompAny(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SeccompAny_Pod_Unconfined(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
 annotations:
   seccomp.security.alpha.kubernetes.io/pod: unconfined
spec:
 containers:
   - name: trustworthy-container
     image: sotrustworthy:latest
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SeccompAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_SeccompAnyMissing_Pod(t *testing.T) {
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

	containers := SeccompAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

// TODO(ajm) more seccomp tests for deployments
