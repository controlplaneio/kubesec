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
		Select("securityContext.capabilities.drop").
		WhereNotNil("securityContext.capabilities.drop")

	if capDrop != nil &&
		capDrop.Count() > 0 &&
		!strings.Contains(fmt.Sprintf("%v", capDrop.Get()), "<nil>") {
		containers += capDrop.Count()
	}

	capDropInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Select("securityContext.capabilities.drop").
		WhereNotNil("securityContext.capabilities.drop")

	if capDropInit != nil &&
		capDropInit.Count() > 0 &&
		!strings.Contains(fmt.Sprintf("%v", capDropInit.Get()), "<nil>") {
		containers += capDropInit.Count()
	}

	return containers
}
