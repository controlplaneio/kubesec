package rules

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/thedevsaddam/gojsonq/v2"
)

// checkSecurityContextFn is the used to check attributes in
// the securityContext of a pod, initContainer or container. It
// expects the json node of a pod or a container.
type checkSecurityContextFn func(*gojsonq.JSONQ) checkSecurityContextResult

// checkSecurityContextResult defines the state of an attribute in security context.
// This helps making decision with attribute precedence.
type checkSecurityContextResult struct {
	// unset defines when an attribute is not set.
	unset bool
	// valid defines when an attribute is correctly configured.
	valid bool
}

func checkSecurityContext(json []byte, checkPodSecurityContext bool, checkFn checkSecurityContextFn) int {
	jq := gojsonq.New().Reader(bytes.NewReader(json))
	spec := getSpecSelector(json)

	// Some attributes can be set in the PodSecurityContext and
	// will be common to all containers but only if the same attribute is
	// not also set in container.securityContext because field values of
	// container.securityContext take precedence over field values of PodSecurityContext.
	var attrValidAtPodLevel bool
	if checkPodSecurityContext {
		jqPod := jq.Copy().From(spec)

		if res := checkFn(jqPod); res.valid {
			attrValidAtPodLevel = true
		}
	}

	containersLen := jq.Copy().From(spec + ".containers").Count()
	initContainersLen := jq.Copy().From(spec + ".initContainers").Count()
	ephemeralContainersLen := jq.Copy().From(spec + ".ephemeralContainers").Count()

	var valid int

	ctnFn := func(containerCount int, containerNode string) {
		for i := 0; i < containerCount; i++ {
			// e.g: .initContainer.[0]
			node := fmt.Sprintf(".%s.[%s]", containerNode, strconv.Itoa(i))
			c := jq.Copy().From(spec + node)

			res := checkFn(c)
			// takes precedence or inherits common setting from pod level.
			if res.valid || (res.unset && attrValidAtPodLevel) {
				valid += 1
			}
		}
	}

	ctnFn(initContainersLen, "initContainers")
	ctnFn(containersLen, "containers")
	ctnFn(ephemeralContainersLen, "ephemeralContainers")

	return valid
}
