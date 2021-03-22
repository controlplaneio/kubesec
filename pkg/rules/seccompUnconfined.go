package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"regexp"
	"strings"
)

// TODO(ajm) this is just an inversion of seccompAny.go and should be refactored to use a shared function
func SeccompUnconfined(json []byte) int {
	containers := 0
	startWordBoundaryRegex := "[\\[ ]"
	endWordBoundaryRegex := "[\\] ]"

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From("metadata.annotations").Get()

	capDropString := fmt.Sprintf("%v", capDrop)

	if capDrop != nil && strings.Contains(capDropString, "seccomp.security.alpha.kubernetes.io/pod:") {
		if strings.Contains(capDropString, "seccomp.security.alpha.kubernetes.io/pod:unconfined") {
			containers++
		}
	} else if capDrop != nil {

		keyNameRegex := "seccomp\\.security\\.alpha\\.kubernetes\\.io/[a-zA-Z-.]+"
		// TODO(ajm) match end of string in regex
		isNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":", capDropString)

		if isNamedPodMatch {
			isUnconfinedNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":unconfined"+endWordBoundaryRegex, capDropString)
			if isUnconfinedNamedPodMatch {
				containers++
			}
		}
	}

	return containers
}
