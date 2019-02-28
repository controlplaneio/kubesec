package rules

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"strings"
)

// TODO(ajm) this should be checking for the name of the pod, not `pod`, so needs wildcard matching
func SeccompAny(json []byte) int {
	containers := 0

	capDrop := gojsonq.New().Reader(bytes.NewReader(json)).
		From("metadata.annotations").Get()

	if capDrop != nil && strings.Contains(fmt.Sprintf("%v", capDrop), "seccomp.security.alpha.kubernetes.io/pod:") {
		// TODO(ajm): tighten these matches, they could be "[seccomp..." or " seccomp...", and "unconfined]" or "unconfined "
		// TODO(ajm): space delimiting matches is insufficient as this could be set to `unconfined blah`
		if !strings.Contains(fmt.Sprintf("%v", capDrop), "seccomp.security.alpha.kubernetes.io/pod:unconfined") {
			containers++
		}
	}

	return containers
}
