package ruler

import (
	"github.com/ghodss/yaml"
	"github.com/in-toto/in-toto-golang/in_toto"
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
  selector:
    matchLabels:
      app: podinfo
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
      labels:
        app: podinfo
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

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)

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

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)

	if len(report.Message) < 1 || !strings.Contains(report.Message, "selector is required") {
		t.Errorf("Got error %v ", report.Message)
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

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)

	// kubeval should error out with:
	// spec.replicas: Invalid type. Expected: integer, given: string
	if len(report.Message) < 1 || !strings.Contains(report.Message, "Invalid type") {
		t.Errorf("Got error %v ", report.Message)
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

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)

	if len(report.Message) < 1 || !strings.Contains(report.Message, "unknown schema") {
		t.Errorf("Got error %v ", report.Message)
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

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)

	if len(report.Message) < 1 || !strings.Contains(report.Message, "not supported") {
		t.Errorf("Got error %v ", report.Message)
	}
}

func TestRuleset_Get_intoto(t *testing.T) {
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

	var reports []Report

	report := NewRuleset(zap.NewNop().Sugar()).generateReport(json)
	reports = append(reports, report)

	link := GenerateInTotoLink(reports, []byte(data)).Signed.(in_toto.Link)

	if len(link.Materials) < 1 || len(link.Products) < 1 {
		t.Errorf("Should have generated a report with at least one material and a product %+v",
			link)
	}
}
