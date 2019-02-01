package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func ReadOnlyRootFilesystem(json []byte) int {
	spec := getSpecSelector(json)

	allContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").Count()

	allInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").Count()

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.readOnlyRootFilesystem", "!=", nil).
		Where("securityContext.readOnlyRootFilesystem", "=", true)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.readOnlyRootFilesystem", "!=", nil).
		Where("securityContext.readOnlyRootFilesystem", "=", true)

	return (allContainers + allInitContainers) - (jqContainers.Count() + jqInitContainers.Count())
}
