package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"strings"
)

func DockerSock(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".volumes").
		Only("hostPath.path")

	if paths != nil && strings.Contains(fmt.Sprintf("%v", paths), "/var/run/docker.sock") {
		found++
	}

	return found
}
