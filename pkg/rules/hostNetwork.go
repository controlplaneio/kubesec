package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func HostNetwork(json []byte) int {
	res := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.hostNetwork").Get()

	if res != nil && res.(bool) {
		return 1
	}

	return 0
}
