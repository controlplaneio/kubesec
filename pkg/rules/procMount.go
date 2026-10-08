package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func ProcMount(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	data := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".volumes").
		Only("hostPath.path")

	paths, ok := data.([]interface{})
	if ok && paths != nil {
		if strings.Contains(fmt.Sprintf("%v", paths), "/proc") {
			found++
		}
	}

	return found
}
