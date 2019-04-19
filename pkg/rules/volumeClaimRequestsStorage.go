package rules

import (
	"bytes"

	"github.com/thedevsaddam/gojsonq"
)

func VolumeClaimRequestsStorage(json []byte) int {
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.volumeClaimTemplates").
		Only("spec.resources.requests.storage")

	if paths != nil {
		found += len(paths.([]interface{}))
	}

	return found
}
