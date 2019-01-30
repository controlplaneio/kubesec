package ruler

import (
	"fmt"
	"github.com/sublimino/kubesec/pkg/rules"
	"go.uber.org/zap"
)

type Ruleset struct {
	Rules  []Rule
	logger *zap.SugaredLogger
}

func NewRuleset(logger *zap.SugaredLogger) *Ruleset {
	list := make([]Rule, 0)

	hostNetworkRule := Rule{
		Predicate: rules.HostNetwork,
		Selector:  ".spec .hostNetwork == true",
		Reason:    "Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostNetworkRule)

	hostPIDRule := Rule{
		Predicate: rules.HostPID,
		Selector:  ".spec .hostPID == true)",
		Reason:    "Sharing the host's PID namespace allows visibility of processes on the host, potentially leaking information such as environment variables and configuration",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostPIDRule)

	readOnlyRootFilesystemRule := Rule{
		Predicate: rules.ReadOnlyRootFilesystem,
		Selector:  "containers[] .securityContext .readOnlyRootFilesystem == true",
		Reason:    "An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    3,
	}
	list = append(list, readOnlyRootFilesystemRule)

	runAsNonRootRule := Rule{
		Predicate: rules.RunAsNonRoot,
		Selector:  "containers[] .securityContext .runAsNonRoot == true",
		Reason:    "Force the running image to run as a non-root user to ensure least privilege",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    10,
	}
	list = append(list, runAsNonRootRule)

	runAsUserRule := Rule{
		Predicate: rules.RunAsUser,
		Selector:  "containers[] .securityContext .runAsUser -gt 10000",
		Reason:    "Run as a high-UID user to avoid conflicts with the host's user table",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    4,
	}
	list = append(list, runAsUserRule)

	privilegedRule := Rule{
		Predicate: rules.Privileged,
		Selector:  "containers[] .securityContext .privileged == true",
		Reason:    "Privileged containers can allow almost completely unrestricted host access",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -30,
	}
	list = append(list, privilegedRule)

	capSysAdminRule := Rule{
		Predicate: rules.CapSysAdmin,
		Selector:  "containers[] .securityContext .capabilities .add == SYS_ADMIN)",
		Reason:    "CAP_SYS_ADMIN is the most privileged capability and should always be avoided",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -30,
	}
	list = append(list, capSysAdminRule)

	dockerSockRule := Rule{
		Predicate: rules.DockerSock,
		Selector:  "volumes[] .hostPath .path == /var/run/docker.sock",
		Reason:    "Mounting the docker.socket leaks information about other containers and can allow container breakout",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, dockerSockRule)

	return &Ruleset{
		Rules:  list,
		logger: logger,
	}
}

func (rs *Ruleset) Run(json []byte) Report {
	report := Report{
		Score: 0,
		Scoring: RuleScoring{
			Advise:   make([]RuleRef, 0),
			Critical: make([]RuleRef, 0),
		},
	}

	var applyedRules int
	for _, rule := range rs.Rules {
		passed, err := rule.Eval(json)

		// skip rule if it doesn't apply to object kind
		switch err.(type) {
		case *NotSupportedError:
			continue
		}

		applyedRules++
		ref := RuleRef{
			Reason:   rule.Reason,
			Selector: rule.Selector,
			Weight:   rule.Weight,
			Link:     rule.Link,
		}

		if passed {
			if rule.Points >= 0 {
				rs.logger.Debugf("positive score rule passed %v", rule.Selector)
				report.Score += rule.Points
			}

			if rule.Points < 0 {
				rs.logger.Debugf("negative score rule passed %v", rule.Selector)
			}
		} else {
			if rule.Points >= 0 {
				rs.logger.Debugf("positive score rule failed %v", rule.Selector)
				report.Scoring.Advise = append(report.Scoring.Advise, ref)
			}

			if rule.Points < 0 {
				rs.logger.Debugf("negative score rule failed %v", rule.Selector)
				report.Scoring.Critical = append(report.Scoring.Critical, ref)
				report.Score += rule.Points
			}
		}
	}

	if applyedRules < 1 {
		report.Error = fmt.Sprintf("This resource kind is not supported")
	} else if report.Score >= 0 {
		report.Success = fmt.Sprintf("Passed with a score of %v points", report.Score)
	} else {
		report.Error = fmt.Sprintf("Failed with a score of %v points", report.Score)
	}

	return report
}
