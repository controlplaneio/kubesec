package rules

import (
	"bytes"

	"github.com/thedevsaddam/gojsonq/v2"
)

// HostUsers checks if the hostUsers field is set to false in the container spec.
// If it is set to false, it returns 1, indicating that the user namespace is being used.
// Otherwise, it returns 0, indicating that the user namespace is not being used.
func HostUsers(json []byte) int {
	spec := getSpecSelector(json)

	res := gojsonq.New().
		Reader(bytes.NewReader(json)).
		From(spec + ".hostUsers").Get()

	// hostUsers: false → Kubernetes creates a separate user‑namespace for the pod, giving
	// the containers their own UID/GID mapping on the node.
	//
	// hostUsers: true (or leaving the field out, which defaults to true), the pod shares
	// the host’s user namespace.

	if res == nil { // if the value is not set, the default is true
		return 0
	}

	// If the value is a boolean, we check its value.
	if v, ok := res.(bool); ok {
		if !v {
			return 1
		}
		return 0
	}
	// default to 0 if the value is not a boolean
	return 0
}
