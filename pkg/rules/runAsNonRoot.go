package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

func RunAsNonRoot(json []byte) int {
	return checkSecurityContext(
		json,
		true,
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.runAsNonRoot").Get()

			v, ok := value.(bool)

			res := checkSecurityContextResult{}
			if !ok {
				res.unset = true
				return res
			}

			if v {
				res.valid = true
			}

			return res
		})
}
