package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RequestsCPU(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("resources.requests.cpu")

	if paths != nil {
		found += len(paths.([]interface{}))
	}

	return found
}
