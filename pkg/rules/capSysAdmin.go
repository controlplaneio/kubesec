package rules

import (
	"fmt"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

func CapSysAdmin(json []byte) int {
	return checkSecurityContext(
		json,
		false, // not present in PodSecurityContext
		func(jq *gojsonq.JSONQ) checkSecurityContextResult {
			value := jq.From("securityContext.capabilities.add").Get()

			res := checkSecurityContextResult{}
			if value == nil {
				res.unset = true
				return res
			}

			if strings.Contains(fmt.Sprintf("%v", value), "SYS_ADMIN") {
				res.valid = true
				return res
			}

			return res
		})
}
