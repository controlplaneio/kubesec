package ruler

import (
	"strings"
	"testing"

	"github.com/in-toto/in-toto-golang/in_toto"
	"go.uber.org/zap"
)

func TestRulesetRun(t *testing.T) {
	tests := []struct {
		name, data string
	}{
		{
			name: "valid yaml",
			data: `
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
`,
		},
		{
			name: "valid json",
			data: `
{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "spec": {
    "selector": {
      "matchLabels": {
        "app": "podinfo"
      }
    },
    "template": {
      "metadata": {
        "annotations": {
          "prometheus.io/scrape": "true"
        },
        "labels": {
          "app": "podinfo"
        }
      },
      "spec": {
        "hostNetwork": true,
        "initContainers": [
          {
            "name": "init1",
            "securityContext": {
              "readOnlyRootFilesystem": true
            }
          },
          {
            "name": "init2",
            "securityContext": {
              "readOnlyRootFilesystem": false
            }
          },
          {
            "name": "init3"
          }
        ],
        "containers": [
          {
            "name": "c1"
          },
          {
            "name": "c2",
            "securityContext": {
              "readOnlyRootFilesystem": false,
              "runAsNonRoot": true,
              "runAsUser": 1001
            }
          },
          {
            "name": "c3",
            "securityContext": {
              "readOnlyRootFilesystem": true
            }
          }
        ]
      }
    }
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultSchemaConfig()
			reports, err := NewRuleset(zap.NewNop().Sugar()).Run("kube.yaml", []byte(tt.data), config)
			if err != nil || len(reports) == 0 {
				t.Fatal(err.Error())
			}

			report := reports[0]

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

		})
	}
}
func TestRulesetRunNoSchemaValidation(t *testing.T) {
	tests := []struct {
		name, data string
	}{
		{
			name: "missing selectors",
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      hostNetwork:
      containers:
        - name: c1
          securityContext:
            readOnlyRootFilesystem: true
            runAsNonRoot: true
        - name: c2
          securityContext:
            readOnlyRootFilesystem: true
            runAsNonRoot: true
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultSchemaConfig()
			config.DisableValidation = true

			reports, err := NewRuleset(zap.NewNop().Sugar()).Run("kube.yaml", []byte(tt.data), config)
			if err != nil || len(reports) == 0 {
				t.Fatal(err.Error())
			}

			report := reports[0]

			if report.Score != 2 {
				t.Errorf("Got score: %d, expected: 2", report.Score)
			}

		})
	}
}

func TestRulesetRunInvalid(t *testing.T) {
	tests := []struct {
		name, data      string
		expectedMessage string
	}{
		{
			name: "missing selectors",
			data: `
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
`,
			expectedMessage: "For field spec: selector is required",
		},
		{
			name: "replicas has wrong type string",
			data: `
---
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: "2"
  selector:
    matchLabels:
      app: podinfo
  template:
    spec:
      initContainers:
        - name: init1
      containers:
        - name: c1
        - name: c2
`,
			expectedMessage: "For field spec.replicas: Invalid type",
		},
		{
			name: "resource kind does not exist",
			data: `
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
`,
			expectedMessage: "could not find schema for Deployment2",
		},
		{
			name: "resource kind is not supported (ConfigMap)",
			data: `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  color: blue
`,
			expectedMessage: "not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultSchemaConfig()
			reports, err := NewRuleset(zap.NewNop().Sugar()).Run("kube.yaml", []byte(tt.data), config)
			if err != nil || len(reports) == 0 {
				t.Fatal(err.Error())
			}

			report := reports[0]

			if report.Message == "" || !strings.Contains(report.Message, tt.expectedMessage) {
				t.Errorf("Got error %v, expected: %v", report.Message, tt.expectedMessage)
			}
		})
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

	config := NewDefaultSchemaConfig()
	reports, err := NewRuleset(zap.NewNop().Sugar()).Run("kube.yaml", []byte(data), config)
	if err != nil || len(reports) == 0 {
		t.Fatal(err.Error())
	}

	link := GenerateInTotoLink(reports, []byte(data)).Signed.(in_toto.Link)

	if len(link.Materials) < 1 || len(link.Products) < 1 {
		t.Errorf("Should have generated a report with at least one material and a product %+v",
			link)
	}
}
