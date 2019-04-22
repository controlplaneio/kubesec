package ruler

import (
	"fmt"
	"github.com/garethr/kubeval/kubeval"
	"github.com/sublimino/kubesec/pkg/rules"
	"go.uber.org/zap"
	"strings"
	"sync"
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
		Selector:  ".spec .hostPID == true",
		Reason:    "Sharing the host's PID namespace allows visibility of processes on the host, potentially leaking information such as environment variables and configuration",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostPIDRule)

	hostIPCRule := Rule{
		Predicate: rules.HostIPC,
		Selector:  ".spec .hostIPC == true",
		Reason:    "Sharing the host's IPC namespace allows container processes to communicate with processes on the host",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostIPCRule)

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
		Selector:  "containers[] .securityContext .capabilities .add == SYS_ADMIN",
		Reason:    "CAP_SYS_ADMIN is the most privileged capability and should always be avoided",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -30,
	}
	list = append(list, capSysAdminRule)

	capDropAnyRule := Rule{
		Predicate: rules.CapDropAny,
		Selector:  "containers[] .securityContext .capabilities .drop",
		Reason:    "Reducing kernel capabilities available to a container limits its attack surface",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, capDropAnyRule)

	capDropAllRule := Rule{
		Predicate: rules.CapDropAll,
		Selector:  "containers[] .securityContext .capabilities .drop | index(\"ALL\")",
		Reason:    "Drop all capabilities and add only those required to reduce syscall attack surface",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, capDropAllRule)

	dockerSockRule := Rule{
		Predicate: rules.DockerSock,
		Selector:  "volumes[] .hostPath .path == /var/run/docker.sock",
		Reason:    "Mounting the docker.socket leaks information about other containers and can allow container breakout",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, dockerSockRule)

	requestsCPURule := Rule{
		Predicate: rules.RequestsCPU,
		Selector:  "containers[] .resources .requests .cpu",
		Reason:    "Enforcing CPU requests aids a fair balancing of resources across the cluster",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, requestsCPURule)

	limitsCPURule := Rule{
		Predicate: rules.LimitsCPU,
		Selector:  "containers[] .resources .limits .cpu",
		Reason:    "Enforcing CPU limits prevents DOS via resource exhaustion",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, limitsCPURule)

	requestsMemoryRule := Rule{
		Predicate: rules.RequestsMemory,
		Selector:  "containers[] .resources .requests .memory",
		Reason:    "Enforcing memory requests aids a fair balancing of resources across the cluster",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, requestsMemoryRule)

	limitsMemoryRule := Rule{
		Predicate: rules.RequestsMemory,
		Selector:  "containers[] .resources .limits .memory",
		Reason:    "Enforcing memory limits prevents DOS via resource exhaustion",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, limitsMemoryRule)

	serviceAccountNameRule := Rule{
		Predicate: rules.ServiceAccountName,
		Selector:  ".spec .serviceAccountName",
		Reason:    "Service accounts restrict Kubernetes API access and should be configured with least privilege",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    3,
	}
	list = append(list, serviceAccountNameRule)

	hostAliasesRule := Rule{
		Predicate: rules.HostAliases,
		Selector:  ".spec .hostAliases",
		Reason:    "Managing /etc/hosts aliases can prevent Docker from modifying the file after a pod's containers have already been started",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    0,
	}
	list = append(list, hostAliasesRule)

	seccompAnyRule := Rule{
		Predicate: rules.SeccompAny,
		Selector:  ".metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/pod\"",
		Reason:    "Seccomp profiles set minimum privilege and secure against unknown threats",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, seccompAnyRule)

	seccompUnconfinedRule := Rule{
		Predicate: rules.SeccompUnconfined,
		Selector:  ".metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/pod\"",
		Reason:    "Unconfined Seccomp profiles have full system call access",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -1,
	}
	list = append(list, seccompUnconfinedRule)

	apparmorAnyRule := Rule{
		Predicate: rules.ApparmorAny,
		Selector:  ".metadata .annotations .\"container.apparmor.security.beta.kubernetes.io/nginx\"",
		Reason:    "Well defined AppArmor policies may provide greater protection from unknown threats. WARNING: NOT PRODUCTION READY",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    3,
	}
	list = append(list, apparmorAnyRule)

	volumeClaimAccessModeReadWriteOnce := Rule{
		Predicate: rules.VolumeClaimAccessModeReadWriteOnce,
		Selector:  ".spec .volumeClaimTemplates[] .spec .accessModes | index(\"ReadWriteOnce\")",
		Reason:    "",
		Kinds:     []string{"StatefulSet"},
		Points:    1,
	}
	list = append(list, volumeClaimAccessModeReadWriteOnce)

	volumeClaimRequestsStorage := Rule{
		Predicate: rules.VolumeClaimRequestsStorage,
		Selector:  ".spec .volumeClaimTemplates[] .spec .resources .requests .storage",
		Reason:    "",
		Kinds:     []string{"StatefulSet"},
		Points:    1,
	}
	list = append(list, volumeClaimRequestsStorage)

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

	// validate resource
	results, err := kubeval.Validate(json, "resource.json")
	if err != nil {
		if strings.Contains(err.Error(), "Problem loading schema from the network") {
      report.Error = "This resource is invalid, unknown schema"
		} else {
			report.Error = err.Error()
		}
		return report
	}
	for _, result := range results {
		if len(result.Errors) > 0 {
			for _, desc := range result.Errors {
				report.Error += desc.String()
			}
		} else if result.Kind == "" {
			report.Error += "This resource is invalid, Kubernetes kind not found"
		}
	}

	if len(report.Error) > 0 {
		return report
	}

	// run rules in parallel
	ch := make(chan RuleRef, len(rs.Rules))
	var wg sync.WaitGroup
	for _, rule := range rs.Rules {
		wg.Add(1)
		go eval(json, rule, ch, &wg)
	}
	wg.Wait()
	close(ch)

	// collect results
	var appliedRules int
	for ruleRef := range ch {
		appliedRules++

		if ruleRef.Containers > 0 {
			if ruleRef.Points >= 0 {
				rs.logger.Debugf("positive score rule matched %v", ruleRef.Selector)
				report.Score += ruleRef.Points
			}

			if ruleRef.Points < 0 {
				rs.logger.Debugf("negative score rule matched %v", ruleRef.Selector)
				report.Score += ruleRef.Points
				report.Scoring.Critical = append(report.Scoring.Critical, ruleRef)
			}
			rs.logger.Debugf("points %v", report.Score)
		} else {
			if ruleRef.Points >= 0 {
				rs.logger.Debugf("positive score rule failed %v", ruleRef.Selector)
				report.Scoring.Advise = append(report.Scoring.Advise, ruleRef)
			}

			if ruleRef.Points < 0 {
				rs.logger.Debugf("negative score rule failed %v", ruleRef.Selector)
				report.Scoring.Critical = append(report.Scoring.Critical, ruleRef)
			}
		}
	}

	if appliedRules < 1 {
		report.Error = fmt.Sprintf("This resource kind is not supported")
	} else if report.Score >= 0 {
		report.Success = fmt.Sprintf("Passed with a score of %v points", report.Score)
	} else {
		report.Error = fmt.Sprintf("Failed with a score of %v points", report.Score)
	}

	return report
}

func eval(json []byte, rule Rule, ch chan RuleRef, wg *sync.WaitGroup) {
	defer wg.Done()

	containers, err := rule.Eval(json)

	// skip rule if it doesn't apply to object kind
	switch err.(type) {
	case *NotSupportedError:
		return
	}

	result := RuleRef{
		Containers: containers,
		Points:     rule.Points,
		Reason:     rule.Reason,
		Selector:   rule.Selector,
		Weight:     rule.Weight,
		Link:       rule.Link,
	}

	ch <- result
}
