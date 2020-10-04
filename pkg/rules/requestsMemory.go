package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RequestsMemory(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("resources.requests.memory")

	if paths != nil {
		found += len(paths.([]interface{}))
	}

	return found
}
