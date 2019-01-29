package ruler

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func hostNetwork(json []byte) int {
	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.containers").
		Where("hostNetwork", "!=", nil).
		Where("hostNetwork", "=", true)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From("spec.template.spec.initContainers").
		Where("hostNetwork", "!=", nil).
		Where("hostNetwork", "=", true)

		//res := jqInitContainers.Get()
		//err := jqInitContainers.Error()
		//fmt.Printf("Error: %v\nResult: %#v\n", err, res)

	return jqContainers.Count() + jqInitContainers.Count()
}
