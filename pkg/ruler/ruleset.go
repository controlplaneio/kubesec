package ruler

import "github.com/sublimino/kubesec/pkg/rules"

type Ruleset struct {
	Rules []Rule
}

func NewRuleset() *Ruleset {
	list := make([]Rule, 0)

	hostNetworkRule := Rule{
		Predicate: rules.HostNetwork,
		Selector:  ".spec .hostNetwork",
		Reason:    "Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter",
		Kinds:     []string{"Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostNetworkRule)

	readOnlyRootFilesystemRule := Rule{
		Predicate: rules.ReadOnlyRootFilesystem,
		Selector:  "containers[] .securityContext .readOnlyRootFilesystem == true",
		Reason:    "An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost",
		Kinds:     []string{"Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    3,
	}
	list = append(list, readOnlyRootFilesystemRule)

	runAsNonRootRule := Rule{
		Predicate: rules.RunAsNonRoot,
		Selector:  "containers[] .securityContext .runAsNonRoot == true",
		Reason:    "Force the running image to run as a non-root user to ensure least privilege",
		Kinds:     []string{"Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    10,
	}
	list = append(list, runAsNonRootRule)

	runAsUserRule := Rule{
		Predicate: rules.RunAsUser,
		Selector:  "containers[] .securityContext .runAsUser -gt 10000",
		Reason:    "Run as a high-UID user to avoid conflicts with the host's user table",
		Kinds:     []string{"Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    4,
	}
	list = append(list, runAsUserRule)

	return &Ruleset{
		Rules: list,
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

	for _, rule := range rs.Rules {
		ok, err := rule.Eval(json)
		switch err.(type) {
		case *NotSupportedError:
			continue
		}
		if !ok {
			ref := RuleRef{
				Reason:   rule.Reason,
				Selector: rule.Selector,
				Weight:   rule.Weight,
				Link:     rule.Link,
			}

			if rule.Points >= 0 {
				report.Scoring.Advise = append(report.Scoring.Advise, ref)
			} else {
				report.Scoring.Critical = append(report.Scoring.Critical, ref)
			}
		}
		report.Score += rule.Points
	}

	return report
}
