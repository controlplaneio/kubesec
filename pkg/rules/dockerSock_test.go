package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_DockerSock_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  volumes:
    - name: docker
      hostPath:
        path: /var/run/docker.sock
    - name: tmp
      hostPath:
        path: /tmp
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := DockerSock(json)
	if containers != 1 {
		t.Errorf("Got %v volumes wanted %v", containers, 1)
	}
}

func Test_DockerSock_DaemonSet(t *testing.T) {
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
        - mountPath: /host/var/run/docker.sock
          name: docker
          readOnly: false
      volumes:
      - name: docker
        hostPath:
         path: /var/run/docker.sock
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := DockerSock(json)
	if containers != 1 {
		t.Errorf("Got %v volumes wanted %v", containers, 1)
	}
}

func Test_DockerSock_Missing(t *testing.T) {
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

	containers := DockerSock(json)
	if containers != 0 {
		t.Errorf("Got %v volumes wanted %v", containers, 0)
	}
}

func Test_DockerSock_NoVolumes(t *testing.T) {
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

	containers := DockerSock(json)
	if containers != 0 {
		t.Errorf("Got %v volumes wanted %v", containers, 0)
	}
}
