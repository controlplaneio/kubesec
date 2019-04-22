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

	if len(report.Error) < 1 || !strings.Contains(report.Error, "unknown schema") {
		t.Errorf("Got error %v ", report.Error)
	}
}

func TestRuleset_Run_not_supported(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  color: blue
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	report := NewRuleset(zap.NewNop().Sugar()).Run(json)

	if len(report.Error) < 1 || !strings.Contains(report.Error, "not supported") {
		t.Errorf("Got error %v ", report.Error)
	}
}
