package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
)

func getSpecSelector(json []byte) string {
	selector := "spec.template.spec"

	jq := gojsonq.New().Reader(bytes.NewReader(json)).From("kind")
	if jq.Error() != nil {
		return selector
	}

	kind := fmt.Sprintf("%s", jq.Get())

	if kind == "Pod" {
		selector = "spec"
	}

	return selector
}
