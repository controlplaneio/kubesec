package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
)

func ServiceAccountName(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".serviceAccountName").
		Where(".serviceAccountName", "!=", nil).
		Where(".serviceAccountName", "!=", "")

	// TODO(ajm) the above `Where` selectors don't do what I'd expect and filter the results
	if fmt.Sprintf("%v", jqContainers.Get()) != "<nil>" {
		return 1
	}

	return 0
}
