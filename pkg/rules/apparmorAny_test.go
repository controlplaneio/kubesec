package rules

import (
	"github.com/ghodss/yaml"
	"testing"
)

func Test_ApparmorAny_Pod(t *testing.T) {
	var data = `
---
# The example Pod utilizing the profile loaded by the sample daemon.

apiVersion: v1
kind: Pod
metadata:
  name: nginx-apparmor
  # Note that the Pod does not need to be in the same namespace as the loader.
  labels:
    app: nginx
  annotations:
    # Tell Kubernetes to apply the AppArmor profile "k8s-nginx".
    # Note that this is ignored if the Kubernetes node is not running version 1.4 or greater.
    container.apparmor.security.beta.kubernetes.io/nginx: localhost/k8s-nginx
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_ApparmorAny_Pod_Unconfined(t *testing.T) {
	var data = `
---
# The example Pod utilizing the profile loaded by the sample daemon.

apiVersion: v1
kind: Pod
metadata:
  name: nginx-apparmor
  # Note that the Pod does not need to be in the same namespace as the loader.
  labels:
    app: nginx
  annotations:
    # Tell Kubernetes to apply the AppArmor profile "k8s-nginx".
    # Note that this is ignored if the Kubernetes node is not running version 1.4 or greater.
    container.apparmor.security.beta.kubernetes.io/nginx: unconfined
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_ApparmorAny_No_Seccomp(t *testing.T) {
	var data = `
---
# The example Pod utilizing the profile loaded by the sample daemon.

apiVersion: v1
kind: Pod
metadata:
  name: nginx-apparmor
  # Note that the Pod does not need to be in the same namespace as the loader.
  labels:
    app: nginx
  annotations:
    # Tell Kubernetes to apply the AppArmor profile "k8s-nginx".
    # Note that this is ignored if the Kubernetes node is not running version 1.4 or greater.
    container.apparmor.security.beta.kubernetes.io/nginx: unconfined
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_ApparmorAny_Named_Pod(t *testing.T) {
	var data = `
---
# The example Pod utilizing the profile loaded by the sample daemon.

apiVersion: v1
kind: Pod
metadata:
  name: nginx-apparmor
  # Note that the Pod does not need to be in the same namespace as the loader.
  labels:
    app: nginx
  annotations:
    # Tell Kubernetes to apply the AppArmor profile "k8s-nginx".
    # Note that this is ignored if the Kubernetes node is not running version 1.4 or greater.
    container.apparmor.security.beta.kubernetes.io/somePodName: localhost/k8s-nginx
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_ApparmorAny_Named_Pod_Special_Chars(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    container.apparmor.security.beta.kubernetes.io/my-Named.Pod: runtime/default
spec:
containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_ApparmorAny_Named_Pod_Special_Chars_Unconfined(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
annotations:
  container.apparmor.security.beta.kubernetes.io/my-Named.Pod: unconfined
spec:
containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_ApparmorAny_Named_Pod_Illegal_Name(t *testing.T) {
	var data = `
---
apiVersion: v1
kind: Pod
metadata:
annotations:
  container.apparmor.security.beta.kubernetes.io/my-Named.Pod(illegal name): runtime/default
spec:
containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := ApparmorAny(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

// TODO(ajm) more apparmor tests for deployments
