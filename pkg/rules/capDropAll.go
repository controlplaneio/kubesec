package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"strings"
)

func CapDropAll(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.drop")

	if capDrop != nil && strings.Contains(strings.ToUpper(fmt.Sprintf("%v", capDrop)), "ALL") {
		containers++
	}

	capDropInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.drop")

	if capDropInit != nil && strings.Contains(strings.ToUpper(fmt.Sprintf("%v", capDropInit)), "ALL") {
		containers++
	}

	return containers
}
