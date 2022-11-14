package rules

import (
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func CapDropAll(json []byte) int {
	return checkSecurityContext(
		json,
		false, // not present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.capabilities.drop").Get()

			res := checkSecurityContextResult{}
			if value == nil {
				res.unset = true
				return res
			}

			if strings.Contains(strings.ToUpper(fmt.Sprintf("%v", value)), "ALL") {
				res.valid = true
			}

			return res
		})
}
