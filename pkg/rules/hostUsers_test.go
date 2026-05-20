package rules

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_HostUsers_Pod(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "Pod with hostUsers set to false",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  hostUsers: false
  containers:
  - name: test-container
    image: test-image:1.14.2
    ports:
    - containerPort: 80
`,
			want: 1,
		},
		{
			name: "Pod with hostUsers set to true",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  hostUsers: true
  containers:
  - name: test-container
    image: test-image:1.14.2
    ports:
    - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Pod with HostUsers set to non-boolean value",
			data: `
---
apiVersion: v1
kind: Pod
metadata:
  name: test-app
spec:
  hostUsers: nonBooleanValue
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
			got := HostUsers(json)
			if got != tc.want {
				t.Errorf("HostUsers() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_HostUsers_Deployment(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "Deployment with Pod Spec hostUsers set to false",
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
      hostUsers: false
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "Deployment with Pod Spec hostUsers set to true",
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
      hostUsers: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Deployment with Pod Spec hostUsers not set",
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
			name: "Deployment withPod Spec hostUsers set to non-boolean value",
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
      hostUsers: nonBooleanValue
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "Deployment with hostUsers set to non-boolean value",
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
      hostUsers: nonBooleanValue
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
			got := HostUsers(json)
			if got != tc.want {
				t.Errorf("HostUsers() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_HostUsers_Daemonset(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "DaemonSet with Pod Spec. hostUsers set to false",
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
      hostUsers: false
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "DaemonSet with Pod Spec. hostUsers set to true",
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
      hostUsers: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "DaemonSet with Pod Spec. hostUsers not set",
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
			name: "DaemonSet with Pod Spec. hostUsers set to non-boolean value",
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
      hostUsers: nonBooleanValue
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
			got := HostUsers(json)
			if got != tc.want {
				t.Errorf("HostUsers() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func Test_HostUsers_StatefulSet(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want int
	}{
		{
			name: "StatefulSet with Pod Spec. hostUsers set to false",
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
      hostUsers: false # user namespace for Pod is enabled
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 1,
		},
		{
			name: "StatefulSet with Pod Spec. hostUsers set to true",
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
      hostUsers: true
      containers:
      - name: test-app
        image: test-app:1.2.3
        ports:
        - containerPort: 80
`,
			want: 0,
		},
		{
			name: "StatefulSet with Pod spec. hostUsers not set",
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
			name: "StatefulSet with Pod spec hostUsers set to non-boolean value",
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
      hostUsers: nonBooleanValue
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
			got := HostUsers(json)
			if got != tc.want {
				t.Errorf("HostUsers() - got %v, wanted %v", got, tc.want)
			}
		})
	}
}
