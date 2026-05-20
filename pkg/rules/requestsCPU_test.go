package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_RequestsCPU_Pod(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: c1
    resources:
      requests:
        cpu: 300m
  - name: c2
    resources:
      requests:
  - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RequestsCPU(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}

	//t.Errorf("FAKE ERROR")

}

// TODO(ajm) will this be validated by kubeval instead?
//func Test_RequestsCPU_Empty_Value_Pod(t *testing.T) {
//  var data = `
//---
//apiVersion: v1
//kind: Pod
//spec:
//  containers:
//  - name: c1
//    resources:
//      requests:
//        cpu:
//  - name: c2
//    resources:
//      requests:
//  - name: c3
//`
//
//  json, err := yaml.YAMLToJSON([]byte(data))
//  if err != nil {
//    t.Fatal(err.Error())
//  }
//
//  containers := RequestsCPU(json)
//  if containers != 0 {
//    t.Errorf("Got %v containers wanted %v", containers, 0)
//  }
//}

func Test_RequestsCPU_With_Limits_Pod(t *testing.T) {
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
     requests:
 - name: c3
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := RequestsCPU(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}
