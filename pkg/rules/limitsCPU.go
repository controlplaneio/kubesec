package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func LimitsCPU(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only(".resources.limits.cpu")

	if paths != nil {
		found++
	}

	return found
}
