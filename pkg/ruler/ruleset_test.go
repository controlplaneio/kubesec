package ruler

import (
	"github.com/ghodss/yaml"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestRuleset_Run(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostNetwork: true
      initContainers:
        - name: init1
          securityContext:
            readOnlyRootFilesystem: true
        - name: init2
          securityContext:
            readOnlyRootFilesystem: false
        - name: init3
      containers:
        - name: c1
        - name: c2
          securityContext:
            readOnlyRootFilesystem: false
            runAsNonRoot: true
            runAsUser: 1001
        - name: c3
          securityContext:
            readOnlyRootFilesystem: true
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	report := NewRuleset(zap.NewNop().Sugar()).Run(json)

	critical := len(report.Scoring.Critical)
	if critical < 1 {
		t.Errorf("Got %v critical rules wanted many", critical)
	}

	advise := len(report.Scoring.Advise)
	if advise < 1 {
		t.Errorf("Got %v advise rules wanted many", advise)
	}

	if report.Score > 0 {
		t.Errorf("Got score %v wanted a negative value", report.Score)
	}
}

func TestRuleset_Run_invalid_network(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostNetwork: 
      initContainers:
        - name: init1
      containers:
        - name: c1
        - name: c2
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	report := NewRuleset(zap.NewNop().Sugar()).Run(json)

	// kubeval should error out with:
	// spec.template.spec.hostNetwork: Invalid type. Expected: boolean, given: null
	if len(report.Error) < 1 || !strings.Contains(report.Error, "Expected: boolean") {
		t.Errorf("Got error %v ", report.Error)
	}
}

func TestRuleset_Run_invalid_replicas(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: "2"
  template:
    spec:
      initContainers:
        - name: init1
      containers:
        - name: c1
        - name: c2
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	report := NewRuleset(zap.NewNop().Sugar()).Run(json)

	// kubeval should error out with:
	// spec.replicas: Invalid type. Expected: integer, given: string
	if len(report.Error) < 1 || !strings.Contains(report.Error, "Expected: integer") {
		t.Errorf("Got error %v ", report.Error)
	}
}

func TestRuleset_Run_invalid_kind(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment2
spec:
  template:
    spec:
      initContainers:
        - name: init1
      containers:
        - name: c1
        - name: c2
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	report := NewRuleset(zap.NewNop().Sugar()).Run(json)

	// kubeval should error out with:
	// Problem loading schema from the network at
	// https://raw.githubusercontent.com/garethr/kubernetes-json-schema/master/master-standalone/deployment2.json:
	// Could not read schema from HTTP, response status is 404 Not Found
	if len(report.Error) < 1 || !strings.Contains(report.Error, "404 Not Found") {
		t.Errorf("Got error %v ", report.Error)
	}
}

//func Test_CapDropAny_Malformed_Fail(t *testing.T) {
//	var data = `
//---
//apiVersion: v1
//kind: Pod
//spec:
//initContainers:
//- name: init1
//containers:
//- name: c1
//  securityContext:
//    capabilities:
//      drop: true
//`
//
//	json, err := yaml.YAMLToJSON([]byte(data))
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	report := NewRuleset(zap.NewNop().Sugar()).Run(json)
//
//	// kubeval should error out with:
//	// spec.replicas: Invalid type. Expected: [array,null], given: boolean
//	if len(report.Error) < 1 || !strings.Contains(report.Error, "Expected: [array,null]") {
//		t.Errorf("Got incorrect error: %v ", report.Error)
//	}
//}

//func Test_CapDropAny_Malformed_Empty_List(t *testing.T) {
//	var data = `
//---
//apiVersion: v1
//kind: Pod
//spec:
//initContainers:
//- name: init1
//containers:
//- name: c1
//  securityContext:
//    capabilities:
//       drop:
//       -
//
//`
//
//	json, err := yaml.YAMLToJSON([]byte(data))
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	report := NewRuleset(zap.NewNop().Sugar()).Run(json)
//
//	fmt.Println("ARSE " + report.Error)
//
//	// kubeval should error out with:
//	// spec.replicas: Invalid type. Expected: integer, given: string
//	if len(report.Error) < 1 || !strings.Contains(report.Error, "Expected: integer") {
//		t.Errorf("Got error %v ", report.Error)
//	}
//}

//func Test_CapDropAny_Malformed_Empty_List_2(t *testing.T) {
//	var data = `
//---
//apiVersion: v1
//kind: Pod
//spec:
// initContainers:
// - name: init1
// containers:
// - name: c1
//   securityContext:
//     capabilities:
//       drop:
//       - THIS_SHOULD_FAIL
//`
//
//	json, err := yaml.YAMLToJSON([]byte(data))
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	report := NewRuleset(zap.NewNop().Sugar()).Run(json)
//
//	fmt.Println("ARSE " + report.Error)
//
//	// kubeval should error out with:
//	// TODO: what should this be? Should we add CAPs to kubeval
//	if len(report.Error) < 1 || !strings.Contains(report.Error, "Expected: SOMETHING VALID") {
//		t.Errorf("Got error %v ", report.Error)
//	}
//}
