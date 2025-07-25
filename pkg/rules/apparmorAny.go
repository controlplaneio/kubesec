package rules

import "github.com/thedevsaddam/gojsonq/v2"

// isApparmorUnconfined checks the appArmorProfile.type field and returns a checkSecurityContextResult struct.
// If the field is set then unset=false. If, on top of that, the value of Unconfined matches the expected value
// then return valid=true.
func isApparmorUnconfined(jq *gojsonq.JSONQ, expectedUnconfined bool) checkSecurityContextResult {
	value := jq.From("securityContext.appArmorProfile.type").Get()

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

func ApparmorAny(json []byte) int {
	return checkSecurityContext(
		json,
		true, // present in Pod Security Context
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			return isApparmorUnconfined(jq, false)
		},
	)
}
