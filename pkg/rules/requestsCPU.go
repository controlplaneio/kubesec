package rules

import (
	"bytes"

	"github.com/thedevsaddam/gojsonq/v2"
)

func RequestsCPU(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	data := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("resources.requests.cpu")

	paths, ok := data.([]interface{})
	if ok && paths != nil {
		found += len(paths)
	}

	return found
}
