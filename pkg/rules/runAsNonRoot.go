package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RunAsNonRoot(json []byte) int {
	spec := getSpecSelector(json)

	allContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").Count()

	allInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").Count()

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsNonRoot", "!=", nil).
		Where("securityContext.runAsNonRoot", "=", true)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsNonRoot", "!=", nil).
		Where("securityContext.runAsNonRoot", "=", true)

	return jqContainers.Count() + jqInitContainers.Count()
}
