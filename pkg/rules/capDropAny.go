package rules

import (
  "bytes"
  "fmt"
  "github.com/thedevsaddam/gojsonq"
  "strings"
)

func CapDropAny(json []byte) int {
	spec := getSpecSelector(json)
	containers := 0

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only("securityContext.capabilities.drop")

	if capDrop != nil && strings.Contains(fmt.Sprintf("%v", capDrop), "[map[drop:[") && ! strings.Contains(fmt.Sprintf("%v", capDrop), "[map[drop:[<nil>]") {
		containers++
	}

	capDropInit := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".initContainers").
		Only("securityContext.capabilities.drop")

  if capDropInit != nil && strings.Contains(fmt.Sprintf("%v", capDropInit), "[map[drop:[") {
		containers++
	}

	return containers
}
