package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"strings"
)

/*
type: container and initContainer
full match in jsonq: no
results per-container, per-spec, per-kind: per-spec (volume)

per-kind: annotations
per-spec: volumes
per-container: security context, image name
*/

func DockerSock(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".volumes").
		Only("hostPath.path")

	if paths != nil &&
		strings.Contains(fmt.Sprintf("%v", paths), "[path:/var/run/docker.sock]") {
		found++
	}

	return found
}
