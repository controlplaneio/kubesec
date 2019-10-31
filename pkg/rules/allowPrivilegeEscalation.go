package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func AllowPrivilegeEscalation(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.allowPrivilegeEscalation", "!=", nil).
		Where("securityContext.allowPrivilegeEscalation", "=", true)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.allowPrivilegeEscalation", "!=", nil).
		Where("securityContext.allowPrivilegeEscalation", "=", true)

	return jqContainers.Count() + jqInitContainers.Count()
}
