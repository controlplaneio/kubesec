package ruler

import (
	"github.com/ghodss/yaml"
	"testing"
)

func TestRule_Eval(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: alpine
          image: alpine
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	rule := &Rule{
		Predicate: hostNetwork,
		Kinds:     []string{"Deployment"},
	}

	ok, err := rule.Eval(json)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !ok {
		t.Errorf("Rule failed when it shouldn't")
	}
}

func TestRule_EvalDoesNotApply(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
kind: StatefulSet
spec:
  template:
    spec:
      containers:
        - name: alpine
          image: alpine
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	rule := &Rule{
		Predicate: hostNetwork,
		Kinds:     []string{"Deployment"},
	}

	_, err = rule.Eval(json)
	if err == nil {
		t.Errorf("Rule succeeded when it shouldn't")
	}
}

func TestRule_EvalNoKind(t *testing.T) {
	var data = `
---
apiVersion: apps/v1
spec:
  template:
    spec:
      containers:
        - name: alpine
          image: alpine
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	rule := &Rule{
		Predicate: hostNetwork,
		Kinds:     []string{"Deployment"},
	}

	_, err = rule.Eval(json)
	if err == nil {
		t.Errorf("Rule succeeded when it shouldn't")
	}
}
