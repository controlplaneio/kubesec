package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func VolumeClaimAccessModeReadWriteOnce(json []byte) int {
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.volumeClaimTemplates").
		Only("spec.accessModes")

	if paths != nil && strings.Contains(fmt.Sprintf("%v", paths), "accessModes:[ReadWriteOnce]") {
		found++
	}

	return found
}
