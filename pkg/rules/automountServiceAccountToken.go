package rules

import (
	"bytes"
	
	"github.com/thedevsaddam/gojsonq/v2"
)

func AutomountServiceAccountToken(json []byte) int {
	spec := getSpecSelector(json)
	
	res := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".automountServiceAccountToken").Get()
	
	if res != nil {
		if v, ok := res.(bool); ok && !v {
			return 1
		}
	}
	
	return 0
}
