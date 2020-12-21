package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq/v2"
)

func HostNetwork(json []byte) int {
	spec := getSpecSelector(json)

	res := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".hostNetwork").Get()

	if res != nil && res.(bool) {
		return 1
	}

	return 0
}
