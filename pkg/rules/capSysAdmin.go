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

func CapSysAdmin(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capAdd := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.add")

	if capAdd != nil && strings.Contains(fmt.Sprintf("%v", capAdd), "SYS_ADMIN") {
		containers += len(capAdd.([]interface{}))
	}

	capAddInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.add")

	if capAddInit != nil && strings.Contains(fmt.Sprintf("%v", capAddInit), "SYS_ADMIN") {
		containers += len(capAddInit.([]interface{}))
	}

	return containers
}
