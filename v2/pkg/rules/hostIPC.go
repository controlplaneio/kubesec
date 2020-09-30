package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func HostIPC(json []byte) int {
	spec := getSpecSelector(json)

	res := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".hostIPC").Get()

	if res != nil && res.(bool) {
		return 1
	}

	return 0
}
