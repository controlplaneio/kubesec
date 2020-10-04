package ruler

import (
	"fmt"
	"testing"

	"github.com/controlplaneio/kubesec/v2/pkg/rules"
	"github.com/ghodss/yaml"
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
          hostNetwork: false
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	rule := &Rule{
		Predicate: rules.HostNetwork,
		Kinds:     []string{"Deployment"},
	}

	matchedContainerCount, err := rule.Eval(json)
	if err != nil {
		t.Fatal(err.Error())
	}
	if matchedContainerCount != 0 {
		t.Errorf(fmt.Sprintf("Rule failed when it shouldn't with count %d", matchedContainerCount))
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
		Predicate: rules.HostNetwork,
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
		Predicate: rules.HostNetwork,
		Kinds:     []string{"Deployment"},
	}

	_, err = rule.Eval(json)
	if err == nil {
		t.Errorf("Rule succeeded when it shouldn't")
	}
}
