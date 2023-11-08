package rules

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

func ReadOnlyRootFilesystem(json []byte) int {
	return checkSecurityContext(
		json,
		false, // not present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.readOnlyRootFilesystem").Get()

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
