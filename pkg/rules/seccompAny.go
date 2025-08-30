package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

// isSeccompUnconfined checks the seccompProfile.type field and returns a checkSecurityContextResult struct.
// If the field is set then unset=false. If, on top of that, the value of Unconfined matches the expected value
// then return valid=true.
func isSeccompUnconfined(jq *gojsonq.JSONQ, expectedUnconfined bool) checkSecurityContextResult {
	value := jq.From("securityContext.seccompProfile.type").Get()

	v, ok := value.(string)

	res := checkSecurityContextResult{}
	if !ok {
		res.unset = true
		return res
	}

	return checkSecurityContextResult{
		valid: (v == "Unconfined") == expectedUnconfined,
	}
}

// SeccompAny retrieves the number of instances in a manifest where the Seccomp profile has been specified
// to a value other than 'Unconfined'
func SeccompAny(json []byte) int {
	return checkSecurityContext(
		json,
		true, // present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			return isSeccompUnconfined(jq, false)
		})
}
