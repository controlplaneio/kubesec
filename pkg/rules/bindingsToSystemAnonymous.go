package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func BindingsToSystemAnonymous(json []byte) int {
	data := gojsonq.New().Reader(bytes.NewReader(json)).
		From("subjects").
		Only("name")

	subjectNames, ok := data.([]interface{})
	if ok && subjectNames != nil {
		if strings.Contains(fmt.Sprintf("%v", subjectNames), "system:anonymous") {
			return 1
		}
	}

	return 0
}
