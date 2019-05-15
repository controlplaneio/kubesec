package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RunAsGroup(json []byte) int {
	spec := getSpecSelector(json)

	jqContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".containers").
		Where("securityContext.runAsGroup", ">", 10000)

	jqInitContainers := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec+".initContainers").
		Where("securityContext.runAsGroup", ">", 10000)

	return jqContainers.Count() + jqInitContainers.Count()
}

// where() ->
// .Count()
// only() ->
// containers += len(capDrop.([]interface{})) - this matches ALL the matches

// to check:
// TODO(ajm): no init containers test
// TODO(ajm): count/len is not being used for iterable resource
// TODO(ajm): reduce complexity
