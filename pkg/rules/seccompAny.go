package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"regexp"
	"strings"
)

// TODO(ajm): tighten these matches, they could be "[seccomp..." or " seccomp...", and "unconfined]" or "unconfined "
// TODO(ajm): space delimiting matches is insufficient as this could be set to `unconfined blah`
func SeccompAny(json []byte) int {
	containers := 0
	startWordBoundaryRegex := "[\\[ ]"
	endWordBoundaryRegex := "[\\] ]"

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From("metadata.annotations").Get()

	capDropString := fmt.Sprintf("%v", capDrop)

	if capDrop != nil && strings.Contains(capDropString, "seccomp.security.alpha.kubernetes.io/pod:") {
		if !strings.Contains(capDropString, "seccomp.security.alpha.kubernetes.io/pod:unconfined") {
			containers++
		}
	} else if capDrop != nil {

		keyNameRegex := "seccomp\\.security\\.alpha\\.kubernetes\\.io/[a-zA-Z-.]+"
		// TODO(ajm) match end of string in regex
		isNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":", capDropString)

		if isNamedPodMatch {
			isUnconfinedNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":unconfined"+endWordBoundaryRegex, capDropString)
			if !isUnconfinedNamedPodMatch {
				containers++
			}
		}
	}

	return containers
}
