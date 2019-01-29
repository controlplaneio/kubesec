package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RunAsUser(json []byte) int {
	allContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.containers").Count()

	allInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.initContainers").Count()

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.containers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsUser", "!=", nil).
		Where("securityContext.runAsUser", ">", 10000)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.initContainers").
		Where("securityContext", "!=", nil).
		Where("securityContext.runAsUser", "!=", nil).
		Where("securityContext.runAsUser", ">", 10000)

		//res := jqInitContainers.Get()
		//err := jqInitContainers.Error()
		//fmt.Printf("Error: %v\nResult: %#v\n", err, res)

	return (allContainers + allInitContainers) - (jqContainers.Count() + jqInitContainers.Count())
}
