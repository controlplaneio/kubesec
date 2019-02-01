package ruler

import (
	"github.com/ghodss/yaml"
	"go.uber.org/zap"
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
