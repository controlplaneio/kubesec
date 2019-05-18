package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
)

/*
type: container and initContainer
full match in jsonq: no
results per-container, per-spec, per-kind: per-spec (volume)

per-kind: annotations
per-spec: volumes
per-container: security context, image name
*/

// TODO(ajm) is this relevant?
func HostAliases(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".hostAliases")

	// TODO(ajm) the above `Where` selectors don't do what I'd expect and filter the results
	if fmt.Sprintf("%v", jqContainers.Get()) != "<nil>" {
		return 1
	}

	return 0
}
