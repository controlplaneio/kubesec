package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func VolumeClaimAccessModeReadWriteOnce(json []byte) int {

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

	data := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.volumeClaimTemplates").
		Only("spec.accessModes")

	paths, ok := data.([]interface{})
	if ok && paths != nil {
		if strings.Contains(fmt.Sprintf("%v", paths), "accessModes:[ReadWriteOnce]") {
			found++
		}
	}

	return found
}
