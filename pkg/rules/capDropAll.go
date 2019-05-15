package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"strings"
)

/*
type: container and initContainer
full match in jsonq: no
results per-container or per-spec: per-container
*/

func CapDropAll(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.drop")

	if capDrop != nil && strings.Contains(fmt.Sprintf("%v", capDrop), "ALL") {
		containers += len(capDrop.([]interface{}))
	}

	capDropInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.drop")

	if capDropInit != nil && strings.Contains(fmt.Sprintf("%v", capDropInit), "ALL") {
		containers += len(capDrop.([]interface{}))
	}

	return containers
}
