package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_VolumeClaimVolumeClaimRequestsStorage(t *testing.T) {
	var data = `
---
apiVersion: "apps/v1"
kind: StatefulSet
spec:
  template:
    spec:
      containers:
      - name: cassandra
        image: gcr.io/google-samples/cassandra:v14
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWrite" ]
  - metadata:
      name: 2nd
    spec:
      accessModes: [ "ReadWrite" ]
      resources:
        requests:
          storage: 1Gi

`
	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := VolumeClaimRequestsStorage(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_VolumeClaimVolumeClaimRequestsStorage_Two_Claims(t *testing.T) {
	var data = `
---
apiVersion: "apps/v1"
kind: StatefulSet
spec:
  template:
    spec:
      containers:
      - name: cassandra
        image: gcr.io/google-samples/cassandra:v14
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWrite" ]
      resources:
        requests:
          storage: 2Gi
  - metadata:
      name: 2nd
    spec:
      accessModes: [ "ReadWrite" ]
      resources:
        requests:
          storage: 1Gi

`
	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := VolumeClaimRequestsStorage(json)
	if containers != 2 {
		t.Errorf("Got %v containers wanted %v", containers, 2)
	}
}

func TestStatefulSetHasNoPVCsAndNoStorageRequests(t *testing.T) {
	var data = `
---
apiVersion: "apps/v1"
kind: StatefulSet
spec:
  template:
    spec:
      containers:
      - name: cassandra
        image: gcr.io/google-samples/cassandra:v14

`
	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := VolumeClaimRequestsStorage(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}

}
