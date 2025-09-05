package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

// SeccompUnconfined retrieves the number of instances in a manifest where the Seccomp profile has been specified
// to a value of 'Unconfined'
func SeccompUnconfined(json []byte) int {
	return checkSecurityContext(
		json,
		true, // present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			return isSeccompUnconfined(jq, true)
		})
}
