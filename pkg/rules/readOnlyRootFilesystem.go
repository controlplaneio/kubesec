package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq/v2"
)

func ReadOnlyRootFilesystem(json []byte) int {
	spec := getSpecSelector(json)

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

	return jqContainers.Count() + jqInitContainers.Count()
}
