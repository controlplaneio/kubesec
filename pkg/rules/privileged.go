package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq/v2"
)

func Privileged(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.privileged", "!=", nil).
		Where("securityContext.privileged", "=", true)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.privileged", "!=", nil).
		Where("securityContext.privileged", "=", true)

	return jqContainers.Count() + jqInitContainers.Count()
}
