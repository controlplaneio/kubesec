package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

func AllowPrivilegeEscalation(json []byte) int {
	return checkSecurityContext(
		json,
		false, // not present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.allowPrivilegeEscalation").Get()

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
