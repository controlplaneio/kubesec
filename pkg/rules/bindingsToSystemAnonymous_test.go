package rules

import (
	"testing"

	"github.com/ghodss/yaml"
)

func Test_BindingsToSystemAnonymous_RoleBinding(t *testing.T) {
	var data = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: anonymous-view-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:anonymous
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	res := BindingsToSystemAnonymous(json)
	if res != 1 {
		t.Errorf("Got %v bindings wanted %v", res, 1)
	}
}

func Test_BindingsToSystemAnonymous_ClusterRoleBinding(t *testing.T) {
	var data = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: anonymous-admin-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:anonymous
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	res := BindingsToSystemAnonymous(json)
	if res != 1 {
		t.Errorf("Got %v bindings wanted %v", res, 1)
	}
}

func Test_BindingsToSystemAnonymous_RoleBinding_No_SystemAnonymous_Binding(t *testing.T) {
	var data = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: anonymous-view-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: ServiceAccount
  name: default
  namespace: kube-system
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	res := BindingsToSystemAnonymous(json)
	if res != 0 {
		t.Errorf("Got %v bindings wanted %v", res, 0)
	}
}

func Test_BindingsToSystemAnonymous_ClusterRoleBinding_No_SystemAnonymous_Binding(t *testing.T) {
	var data = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: anonymous-admin-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: ServiceAccount
  name: default
  namespace: kube-system
`

	json, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		t.Fatal(err.Error())
	}

	res := BindingsToSystemAnonymous(json)
	if res != 0 {
		t.Errorf("Got %v bindings wanted %v", res, 0)
	}
}
