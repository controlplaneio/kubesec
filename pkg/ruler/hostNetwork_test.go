package ruler

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_hostNetworkEnabled(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: hostnetwork1
        - name: hostnetwork2
          hostNetwork: false
        - name: hostnetwork3
          hostNetwork: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_hostNetworkManyEnabled(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: hostnetwork1
          hostNetwork: true
        - name: hostnetwork2
          hostNetwork: false
        - name: hostnetwork3
          hostNetwork: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func Test_hostNetworkInitContainersManyEnabled(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      initContainers:
        - name: init1
          hostNetwork: true
      containers:
        - name: hostnetwork1
          hostNetwork: true
        - name: hostnetwork2
          hostNetwork: false
        - name: hostnetwork3
          hostNetwork: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 3 {
		t.Errorf("Got %v containers wanted %v", containers, 3)
	}
}

func Test_hostNetworkDisabled(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: hostnetwork1
        - name: hostnetwork2
          hostNetwork: false
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_hostNetworkMissing(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: hostnetwork1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_hostNetworkNoContainers(t *testing.T) {

	var data = `
---
apiVersion: extensions/v1beta1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: kubesec
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := hostNetwork(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
