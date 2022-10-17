package rules

import (
	"bytes"

	"github.com/thedevsaddam/gojsonq/v2"
)

func VolumeClaimRequestsStorage(json []byte) int {
	volumeClaims := 0
	// count all volumeClaimTemplates
	volumeClaimTemplates := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.volumeClaimTemplates").
		Only("spec")

	volumeClaims += len(volumeClaimTemplates.([]interface{}))

	// pass test if no PVCs are included in statefulset (which is legal)
	if volumeClaims == 0 {
		return 1
	}

	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.volumeClaimTemplates").
		Only("spec.resources.requests.storage")

	found += len(paths.([]interface{}))

	return found
}
