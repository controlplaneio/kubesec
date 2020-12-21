package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"strings"
)

func CapSysAdmin(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capAdd := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.add")

	if capAdd != nil && strings.Contains(fmt.Sprintf("%v", capAdd), "SYS_ADMIN") {
		containers++
	}

	capAddInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.add")

	if capAddInit != nil && strings.Contains(fmt.Sprintf("%v", capAddInit), "SYS_ADMIN") {
		containers++
	}

	return containers
}
