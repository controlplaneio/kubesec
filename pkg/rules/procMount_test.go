package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_ProcMount_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  volumes:
    - name: proc
      hostPath:
        path: /proc
    - name: tmp
      hostPath:
        path: /tmp
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ProcMount(json)
	if containers != 1 {
		t.Errorf("Got %v volumes wanted %v", containers, 1)
	}
}

func Test_ProcMount_DaemonSet(t *testing.T) {
	var data = `
---
apiVersion: extensions/v1beta1
kind: DaemonSet
spec:
  template:
    spec:
      containers:
      - name: c1
        volumeMounts:
        - mountPath: /tmp
          name: proc
          readOnly: false
      volumes:
      - name: proc
        hostPath:
         path: /proc
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ProcMount(json)
	if containers != 1 {
		t.Errorf("Got %v volumes wanted %v", containers, 1)
	}
}

func Test_ProcMount_Missing(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  volumes:
    - name: tmp
      hostPath:
        path: /tmp
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ProcMount(json)
	if containers != 0 {
		t.Errorf("Got %v volumes wanted %v", containers, 0)
	}
}

func Test_ProcMount_NoVolumes(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostNetwork: false
      containers:
        - name: c1
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ProcMount(json)
	if containers != 0 {
		t.Errorf("Got %v volumes wanted %v", containers, 0)
	}
}
