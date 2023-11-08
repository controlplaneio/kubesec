package rules

import (
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func CapDropAny(json []byte) int {
	return checkSecurityContext(
		json,
		false, // not present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.capabilities.drop").Get()

			v, ok := value.([]interface{})

			res := checkSecurityContextResult{}
			if !ok {
				res.unset = true
				return res
			}

			if len(v) > 0 &&
				!strings.Contains(fmt.Sprintf("%v", v), "<nil>") {
				res.valid = true
			}

			return res
		})
}
