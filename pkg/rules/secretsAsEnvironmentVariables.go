package rules

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/thedevsaddam/gojsonq/v2"
)

func SecretsAsEnvironmentVariables(json []byte) int {
	spec := getSpecSelector(json)
	jq := gojsonq.New().Reader(bytes.NewReader(json))
	found := 0

	containerPaths := []string{
		spec + ".containers",
		spec + ".initContainers",
		spec + ".ephemeralContainers",
	}

	// Check for secrets in env, e.g. .containers.[0].env.[0].valueFrom.secretKeyRef
	checkEnv := func(jq *gojsonq.JSONQ, containerNode string) bool {
		envNode := containerNode + ".env"
		envCount := jq.Copy().From(envNode).Count()
		for i := 0; i < envCount; i++ {
			secretKeyRef := jq.Copy().From(fmt.Sprintf("%s.[%s].%s", envNode, strconv.Itoa(i), "valueFrom.secretKeyRef")).Get()
			if secretKeyRef != nil {
				return true
			}
		}
		return false
	}

	// Check for secrets in envFrom, e.g. .containers.[0].envFrom.[0].secretRef
	checkEnvFrom := func(jq *gojsonq.JSONQ, containerNode string) bool {
		envNode := containerNode + ".envFrom"
		envCount := jq.Copy().From(envNode).Count()
		for i := 0; i < envCount; i++ {
			secretRef := jq.Copy().From(fmt.Sprintf("%s.[%s].%s", envNode, strconv.Itoa(i), "secretRef")).Get()
			if secretRef != nil {
				return true
			}
		}
		return false
	}

	for _, containerPath := range containerPaths {
		containerCount := jq.Copy().From(containerPath).Count()

		for i := 0; i < containerCount; i++ {
			containerNode := fmt.Sprintf("%s.[%s]", containerPath, strconv.Itoa(i))

			hasSecretsinEnv := checkEnv(jq, containerNode)
			hasSecretsinEnvFrom := checkEnvFrom(jq, containerNode)

			if hasSecretsinEnv || hasSecretsinEnvFrom {
				found++
			}
		}
	}

	return found
}
