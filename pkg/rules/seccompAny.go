package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

// SeccompAny retrieves the number of instances in a manifest where the Seccomp profile has been specified
// to a value other than 'Unconfined'
func SeccompAny(json []byte) int {
	return checkSecurityContext(
		json,
		true, // present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.seccompProfile.type").Get()

			v, ok := value.(string)

			res := checkSecurityContextResult{}
			if !ok {
				res.unset = true
				return res
			}

			if v != "Unconfined" {
				res.valid = true
			}

			return res
		})
}
