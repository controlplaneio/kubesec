package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func HostPID(json []byte) int {
	spec := getSpecSelector(json)

	res := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".hostPID").Get()

	if res != nil && res.(bool) {
		return 1
	}

	return 0
}
