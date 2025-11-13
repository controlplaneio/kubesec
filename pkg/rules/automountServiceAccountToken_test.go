package rules

import (
	"testing"

	"github.com/ghodss/yaml"
)

func Test_AutomountServiceAccountToken_Pod(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "Pod with automountServiceAccountToken set to false",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  automountServiceAccountToken: false
  containers:
  - name: test-container
    image: test-image:1.14.2
    ports:
    - containerPort: 80
`,
			want: 1,
		},
		{
			name: "Pod with automountServiceAccountToken set to true",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  automountServiceAccountToken: true
  containers:
  - name: test-container
    image: test-image:1.14.2
    ports:
    - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Pod with automountServiceAccountToken set to non-boolean value",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  automountServiceAccountToken: nonBooleanValue
  containers:
  - name: test-container
    image: test-image:1.2.3
    ports:
    - containerPort: 80
`,
			want: 0,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Convert YAML to JSON
			json, err := yaml.YAMLToJSON([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			got := AutomountServiceAccountToken(json)
			if got != tc.want {
				t.Errorf("AutomountServiceAccountToken() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_AutomountServiceAccountToken_Deployment(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "Deployment with automountServiceAccountToken set to false",
			data: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: false
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "Deployment with automountServiceAccountToken set to true",
			data: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Deployment with automountServiceAccountToken not set",
			data: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Deployment with automountServiceAccountToken set to non-boolean value",
			data: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: nonBooleanValue
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Deployment with automountServiceAccountToken set to non-boolean value",
			data: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: nonBooleanValue
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Convert YAML to JSON
			json, err := yaml.YAMLToJSON([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			got := AutomountServiceAccountToken(json)
			if got != tc.want {
				t.Errorf("AutomountServiceAccountToken() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_AutomountServiceAccountToken_Daemonset(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "DaemonSet with automountServiceAccountToken set to false",
			data: `
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test-app-daemonset
spec:
  selector:
    matchLabels:
      name: test-app
  template:
    metadata:
      labels:
        name: test-app
    spec:
      automountServiceAccountToken: false
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "DaemonSet with automountServiceAccountToken set to true",
			data: `
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test-app-daemonset
spec:
  selector:
    matchLabels:
      name: test-app
  template:
    metadata:
      labels:
        name: test-app
    spec:
      automountServiceAccountToken: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "DaemonSet with automountServiceAccountToken not set",
			data: `
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test-app-daemonset
spec:
  selector:
    matchLabels:
      name: test-app
  template:
    metadata:
      labels:
        name: test-app
    spec:
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "DaemonSet with automountServiceAccountToken set to non-boolean value",
			data: `
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test-app-daemonset
spec:
  selector:
    matchLabels:
      name: test-app
  template:
    metadata:
      labels:
        name: test-app
    spec:
      automountServiceAccountToken: nonBooleanValue
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Convert YAML to JSON
			json, err := yaml.YAMLToJSON([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			got := AutomountServiceAccountToken(json)
			if got != tc.want {
				t.Errorf("AutomountServiceAccountToken() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_AutomountServiceAccountToken_StatefulSet(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "StatefulSet with automountServiceAccountToken set to false",
			data: `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-app-statefulset
spec:
  serviceName: "test-app-service"
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: false
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "StatefulSet with automountServiceAccountToken set to true",
			data: `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-app-statefulset
spec:
  serviceName: "test-app-service"
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "StatefulSet with automountServiceAccountToken not set",
			data: `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-app-statefulset
spec:
  serviceName: "test-app-service"
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "StatefulSet with automountServiceAccountToken set to non-boolean value",
			data: `
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-app-statefulset
spec:
  serviceName: "test-app-service"
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      automountServiceAccountToken: nonBooleanValue
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Convert YAML to JSON
			json, err := yaml.YAMLToJSON([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			got := AutomountServiceAccountToken(json)
			if got != tc.want {
				t.Errorf("AutomountServiceAccountToken() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}
