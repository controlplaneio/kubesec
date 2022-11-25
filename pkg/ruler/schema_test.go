package ruler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

// TestValidateSchemaHTTP is a quick test to ensure the schemas
// are fetched from the right location which depends of the GVK and
// if a Kubernetes version is specified. This does not intend to cover
// every possible location as this is already tested in kubeconform,
// just ensuring using SchemaConfig works as expected to pass parameters.
func TestValidateSchemaHTTP(t *testing.T) {
	var tests = []struct {
		name               string
		data               string
		kubernetesVersion  string
		expectedSchemaPath string
		disableValidation  bool
	}{
		{
			name: "kubernetes version not specified",
			data: `
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: test
  name: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - image: test
        name: test
`,
			kubernetesVersion:  "",
			expectedSchemaPath: "/master-standalone-strict/deployment-apps-v1.json",
		},
		{
			name: "kubernetes version specified",
			data: `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    app: test
  name: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - image: test
        name: test
`,
			kubernetesVersion:  "1.25.4",
			expectedSchemaPath: "/v1.25.4-standalone-strict/pod-v1.json",
		},
		{
			name: "validation disabled",
			data: `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    app: test
  name: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - image: test
        name: test
`,
			disableValidation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			called := false
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.URL.Path != tt.expectedSchemaPath {
					t.Errorf("Path: expected %s, actual %s", tt.expectedSchemaPath, r.URL.Path)
				}
			})

			config := NewDefaultSchemaConfig()
			config.DisableValidation = tt.disableValidation
			config.Locations = []string{server.URL}
			config.ValidatorOpts.KubernetesVersion = tt.kubernetesVersion

			reports, err := NewRuleset(zap.NewNop().Sugar()).Run("kube.yaml", []byte(tt.data), config)
			if err != nil || len(reports) == 0 {
				t.Fatal(err.Error())
			}

			if tt.disableValidation && called {
				t.Error("Validation is disabled but validator tried to fetch a schema")
			}

			// We don't actually test validation with a proper schema
			// as this is already done in the ruleset checks.
		})
	}
}
