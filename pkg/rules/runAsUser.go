package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

func RunAsUser(json []byte) int {
	return checkSecurityContext(
		json,
		true,
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.runAsUser").Get()

			v, ok := value.(float64)

			res := checkSecurityContextResult{}
			if !ok {
				res.unset = true
				return res
			}

			if v > 10000 {
				res.valid = true
			}

			return res
		})
}
