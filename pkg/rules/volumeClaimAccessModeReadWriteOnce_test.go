package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_VolumeClaimAccessModeReadWriteOnce(t *testing.T) {
	var data = `
---
apiVersion: "apps/v1"
kind: StatefulSet
metadata:
  name: cassandra
  labels:
     app: cassandra
spec:
  serviceName: cassandra
  replicas: 3
  selector:
    matchLabels:
      app: cassandra
  template:
    spec:
      containers:
      - name: cassandra
        image: gcr.io/google-samples/cassandra:v14
        securityContext:
          capabilities:
            add:
              - IPC_LOCK
        volumeMounts:
        - name: cassandra-data
          mountPath: /cassandra_data
        resources:
          requests:
           memory: 100Mi

  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
  - metadata:
      name: 2nd
    spec:
      accessModes: [ "ReadWrite" ]

`
	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := VolumeClaimAccessModeReadWriteOnce(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func TestStatefulSetHasNoPVCsAndNoAccessModes(t *testing.T) {
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

	containers := VolumeClaimAccessModeReadWriteOnce(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}

}
