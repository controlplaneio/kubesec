package rules

import "github.com/thedevsaddam/gojsonq/v2"

func ApparmorUnconfined(json []byte) int {
	return checkSecurityContext(
		json,
		true, // can be found in Pod Security Context
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			return isApparmorUnconfined(jq, true)
		},
	)
}
