package rules

import (
	"testing"

	"github.com/ghodss/yaml"
)

func Test_SecretsAsEnvironmentVariables_Pod_Env(t *testing.T) {

	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: env-secret
spec:
  containers:
  - name: envars-test-container
    image: nginx
    env:
    - name: SECRET_USERNAME
      valueFrom:
        secretKeyRef:
          name: backend-user
          key: backend-username
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SecretsAsEnvironmentVariables_Pod_EnvFrom(t *testing.T) {

	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: envfrom-secret
spec:
  containers:
  - name: envars-test-container
    image: nginx
    envFrom:
    - secretRef:
        name: test-secret
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SecretsAsEnvironmentVariables_Deployment_Env(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: env-secret
spec:
  template:
    spec:
      containers:
      - name: envars-test-container
        image: nginx
        env:
        - name: SECRET_USERNAME
          valueFrom:
            secretKeyRef:
              name: backend-user
              key: backend-username
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SecretsAsEnvironmentVariables_Deployment_EnvFrom(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: envfrom-secret
spec:
  template:
    spec:
      containers:
      - name: envars-test-container
        image: nginx
        envFrom:
        - secretRef:
            name: test-secret
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SecretsAsEnvironmentVariables_CronJob_Env(t *testing.T) {

	var data = `
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: env-secret
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: envars-test-container
            image: nginx
            env:
            - name: SECRET_USERNAME
              valueFrom:
                secretKeyRef:
                  name: backend-user
                  key: backend-username
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 1 {
		t.Errorf("Got %v containers wanted %v", containers, 1)
	}
}

func Test_SecretsAsEnvironmentVariables_Pod_No_Env_Secrets(t *testing.T) {

	var data = `
---
apiVersion: v1
kind: Pod
metadata:
  name: no-env-secrets
spec:
  containers:
  - name: envars-test-container
    image: nginx
    env:
    - name: SECRET_USERNAME
      value: bob
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}

func Test_SecretsAsEnvironmentVariables_Deployment_No_Env_Secrets(t *testing.T) {

	var data = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: no-env-secrets
spec:
  template:
    spec:
      containers:
      - name: envars-test-container
        image: nginx
        env:
        - name: SECRET_USERNAME
          value: bob
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	containers := SecretsAsEnvironmentVariables(json)
	if containers != 0 {
		t.Errorf("Got %v containers wanted %v", containers, 0)
	}
}
