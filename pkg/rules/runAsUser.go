package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RunAsUser(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsUser", "!=", nil).
		Where("securityContext.runAsUser", ">", 10000)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsUser", "!=", nil).
		Where("securityContext.runAsUser", ">", 10000)

	return jqContainers.Count() + jqInitContainers.Count()
}
