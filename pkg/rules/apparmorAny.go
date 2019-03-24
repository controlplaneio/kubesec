package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"regexp"
	"strings"
)

// TODO(ajm): tighten these matches, they could be "[apparmor..." or " apparmor...", and "unconfined]" or "unconfined "
// TODO(ajm): space delimiting matches is insufficient as this could be set to `unconfined blah`
func ApparmorAny(json []byte) int {
	containers := 0
	startWordBoundaryRegex := "[\\[ ]"
	endWordBoundaryRegex := "[\\] ]"

	annotations := gojsonq.New().Reader(bytes.NewReader(json)).
		From("metadata.annotations").Get()

	annotationsString := fmt.Sprintf("%v", annotations)

	if annotations != nil && strings.Contains(annotationsString, "container.apparmor.security.beta.kubernetes.io/pod:") {
		if !strings.Contains(annotationsString, "container.apparmor.security.beta.kubernetes.io/pod:unconfined") {
			containers++
		}
	} else if annotations != nil {

		keyNameRegex := "container\\.apparmor\\.security\\.beta\\.kubernetes\\.io/[a-zA-Z-.]+"
		// TODO(ajm) match end of string in regex
		isNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":", annotationsString)

		if isNamedPodMatch == true {
			isUnconfinedNamedPodMatch, _ := regexp.MatchString(startWordBoundaryRegex+keyNameRegex+":unconfined"+endWordBoundaryRegex, annotationsString)
			if !isUnconfinedNamedPodMatch {
				containers++
			}
		}
	}

	return containers
}
