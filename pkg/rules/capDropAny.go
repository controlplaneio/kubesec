package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"strings"
)

func CapDropAny(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.drop")

	if capDrop != nil &&
		len(capDrop.([]interface{})) > 0 &&
		!strings.Contains(fmt.Sprintf("%v", capDrop), "<nil>") {
		containers++
	}

	capDropInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.drop")

	if capDropInit != nil &&
		len(capDropInit.([]interface{})) > 0 &&
		!strings.Contains(fmt.Sprintf("%v", capDropInit), "<nil>") {
		containers++
	}

	return containers
}
