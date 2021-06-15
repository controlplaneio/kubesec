package ruler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/controlplaneio/kubesec/v2/pkg/rules"
	"github.com/ghodss/yaml"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/instrumenta/kubeval/kubeval"
	"github.com/thedevsaddam/gojsonq/v2"
	"go.uber.org/zap"
)

type Ruleset struct {
	Rules  []Rule
	logger *zap.SugaredLogger
}

type InvalidInputError struct {
}

func (e *InvalidInputError) Error() string {
	return "Invalid input"
}

func NewRuleset(logger *zap.SugaredLogger) *Ruleset {
	list := make([]Rule, 0)

	hostNetworkRule := Rule{
		Predicate: rules.HostNetwork,
		ID:        "HostNetwork",
		Selector:  ".spec .hostNetwork == true",
		Reason:    "Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostNetworkRule)

	hostPIDRule := Rule{
		Predicate: rules.HostPID,
		ID:        "HostPID",
		Selector:  ".spec .hostPID == true",
		Reason:    "Sharing the host's PID namespace allows visibility of processes on the host, potentially leaking information such as environment variables and configuration",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostPIDRule)

	hostIPCRule := Rule{
		Predicate: rules.HostIPC,
		ID:        "HostIPC",
		Selector:  ".spec .hostIPC == true",
		Reason:    "Sharing the host's IPC namespace allows container processes to communicate with processes on the host",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, hostIPCRule)

	readOnlyRootFilesystemRule := Rule{
		Predicate: rules.ReadOnlyRootFilesystem,
		ID:        "ReadOnlyRootFilesystem",
		Selector:  "containers[] .securityContext .readOnlyRootFilesystem == true",
		Reason:    "An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    3,
	}
	list = append(list, readOnlyRootFilesystemRule)

	runAsNonRootRule := Rule{
		Predicate: rules.RunAsNonRoot,
		ID:        "RunAsNonRoot",
		Selector:  "containers[] .securityContext .runAsNonRoot == true",
		Reason:    "Force the running image to run as a non-root user to ensure least privilege",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    10,
	}
	list = append(list, runAsNonRootRule)

	runAsUserRule := Rule{
		Predicate: rules.RunAsUser,
		ID:        "RunAsUser",
		Selector:  "containers[] .securityContext .runAsUser -gt 10000",
		Reason:    "Run as a high-UID user to avoid conflicts with the host's user table",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
		Advise:    4,
	}
	list = append(list, runAsUserRule)

	privilegedRule := Rule{
		Predicate: rules.Privileged,
		ID:        "Privileged",
		Selector:  "containers[] .securityContext .privileged == true",
		Reason:    "Privileged containers can allow almost completely unrestricted host access",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -30,
	}
	list = append(list, privilegedRule)

	capSysAdminRule := Rule{
		Predicate: rules.CapSysAdmin,
		ID:        "CapSysAdmin",
		Selector:  "containers[] .securityContext .capabilities .add == SYS_ADMIN",
		Reason:    "CAP_SYS_ADMIN is the most privileged capability and should always be avoided",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -30,
	}
	list = append(list, capSysAdminRule)

	capDropAnyRule := Rule{
		Predicate: rules.CapDropAny,
		ID:        "CapDropAny",
		Selector:  "containers[] .securityContext .capabilities .drop",
		Reason:    "Reducing kernel capabilities available to a container limits its attack surface",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, capDropAnyRule)

	capDropAllRule := Rule{
		Predicate: rules.CapDropAll,
		ID:        "CapDropAll",
		Selector:  "containers[] .securityContext .capabilities .drop | index(\"ALL\")",
		Reason:    "Drop all capabilities and add only those required to reduce syscall attack surface",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, capDropAllRule)

	dockerSockRule := Rule{
		Predicate: rules.DockerSock,
		ID:        "DockerSock",
		Selector:  "volumes[] .hostPath .path == /var/run/docker.sock",
		Reason:    "Mounting the docker.socket leaks information about other containers and can allow container breakout",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -9,
	}
	list = append(list, dockerSockRule)

	requestsCPURule := Rule{
		Predicate: rules.RequestsCPU,
		ID:        "RequestsCPU",
		Selector:  "containers[] .resources .requests .cpu",
		Reason:    "Enforcing CPU requests aids a fair balancing of resources across the cluster",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, requestsCPURule)

	limitsCPURule := Rule{
		Predicate: rules.LimitsCPU,
		ID:        "LimitsCPU",
		Selector:  "containers[] .resources .limits .cpu",
		Reason:    "Enforcing CPU limits prevents DOS via resource exhaustion",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, limitsCPURule)

	requestsMemoryRule := Rule{
		Predicate: rules.RequestsMemory,
		ID:        "RequestsMemory",
		Selector:  "containers[] .resources .requests .memory",
		Reason:    "Enforcing memory requests aids a fair balancing of resources across the cluster",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, requestsMemoryRule)

	limitsMemoryRule := Rule{
		Predicate: rules.LimitsMemory,
		ID:        "LimitsMemory",
		Selector:  "containers[] .resources .limits .memory",
		Reason:    "Enforcing memory limits prevents DOS via resource exhaustion",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, limitsMemoryRule)

	serviceAccountNameRule := Rule{
		Predicate: rules.ServiceAccountName,
		ID:        "ServiceAccountName",
		Selector:  ".spec .serviceAccountName",
		Reason:    "Service accounts restrict Kubernetes API access and should be configured with least privilege",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    3,
	}
	list = append(list, serviceAccountNameRule)

	hostAliasesRule := Rule{
		Predicate: rules.HostAliases,
		ID:        "HostAliases",
		Selector:  ".spec .hostAliases",
		Reason:    "Managing /etc/hosts aliases can prevent the container from modifying the file after a pod's containers have already been started. DNS should be managed by the orchestrator",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -3,
	}
	list = append(list, hostAliasesRule)

	seccompAnyRule := Rule{
		Predicate: rules.SeccompAny,
		ID:        "SeccompAny",
		Selector:  ".metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/pod\"",
		Reason:    "Seccomp profiles set minimum privilege and secure against unknown threats",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    1,
	}
	list = append(list, seccompAnyRule)

	seccompUnconfinedRule := Rule{
		Predicate: rules.SeccompUnconfined,
		ID:        "SeccompUnconfined",
		Selector:  ".metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/pod\"",
		Reason:    "Unconfined Seccomp profiles have full system call access",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -1,
	}
	list = append(list, seccompUnconfinedRule)

	apparmorAnyRule := Rule{
		Predicate: rules.ApparmorAny,
		ID:        "ApparmorAny",
		Selector:  ".metadata .annotations .\"container.apparmor.security.beta.kubernetes.io/nginx\"",
		Reason:    "Well defined AppArmor policies may provide greater protection from unknown threats. WARNING: NOT PRODUCTION READY",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    3,
	}
	list = append(list, apparmorAnyRule)

	volumeClaimAccessModeReadWriteOnce := Rule{
		Predicate: rules.VolumeClaimAccessModeReadWriteOnce,
		ID:        "VolumeClaimAccessModeReadWriteOnce",
		Selector:  ".spec .volumeClaimTemplates[] .spec .accessModes | index(\"ReadWriteOnce\")",
		Reason:    "",
		Kinds:     []string{"StatefulSet"},
		Points:    1,
	}
	list = append(list, volumeClaimAccessModeReadWriteOnce)

	volumeClaimRequestsStorage := Rule{
		Predicate: rules.VolumeClaimRequestsStorage,
		ID:        "VolumeClaimRequestsStorage",
		Selector:  ".spec .volumeClaimTemplates[] .spec .resources .requests .storage",
		Reason:    "",
		Kinds:     []string{"StatefulSet"},
		Points:    1,
	}
	list = append(list, volumeClaimRequestsStorage)

	allowPrivilegeEscalation := Rule{
		Predicate: rules.AllowPrivilegeEscalation,
		ID:        "AllowPrivilegeEscalation",
		Selector:  "containers[] .securityContext .allowPrivilegeEscalation == true",
		Reason:    "",
		Kinds:     []string{"Pod", "Deployment", "StatefulSet", "DaemonSet"},
		Points:    -7,
	}
	list = append(list, allowPrivilegeEscalation)

	return &Ruleset{
		Rules:  list,
		logger: logger,
	}
}

func (rs *Ruleset) Run(fileName string, fileBytes []byte, schemaDir string) ([]Report, error) {
	reports := make([]Report, 0)

	isJSON := json.Valid(fileBytes)
	if isJSON {
		report := rs.generateReport(fileName, fileBytes, schemaDir)
		reports = append(reports, report)
	} else {
		lineBreak := detectLineBreak(fileBytes)
		bits := bytes.Split(fileBytes, []byte(lineBreak+"---"+lineBreak))
		for i, d := range bits {
			doc := bytes.TrimSpace(d)

			// If empty or just a header
			if len(doc) == 0 || (len(doc) == 3 && string(doc) == "---") {
				// if we're at the end and there are no reports
				if len(bits) == i+1 && len(reports) == 0 {
					rs.logger.Debugf("empty and no records, erroring")
					return nil, &InvalidInputError{}
				}
				rs.logger.Debugf("empty but still more docs, continuing")
				continue
			}
			data, err := yaml.YAMLToJSON(doc)
			if err != nil {
				return reports, err
			}
			report := rs.generateReport(fileName, data, schemaDir)
			reports = append(reports, report)
		}
	}

	return reports, nil
}

func GenerateInTotoLink(reports []Report, fileBytes []byte) in_toto.Metablock {

	var linkMb in_toto.Metablock

	materials := make(map[string]interface{})
	request := make(map[string]interface{})

	// INFO: it appears that the last newline of the yaml is removed when
	// receiving, which makes the integrity check fail on other implementations
	fileBytes = append(fileBytes, 10)

	request["sha256"] = fmt.Sprintf("%x", sha256.Sum256([]uint8(fileBytes)))

	// TODO: the filename should be a parameter passed to the report (as it is
	// very likely other filenames will exist in supply chains)
	materials["deployment.yml"] = request

	products := make(map[string]interface{})
	for _, report := range reports {
		reportArtifact := make(map[string]interface{})
		// FIXME: encoding as json now for integrity check, this is the wrong way
		// to compute the hash over the result. Also, some error checking would be
		// more than ideal.
		reportValue, _ := json.Marshal(report)
		reportArtifact["sha256"] =
			fmt.Sprintf("%x", sha256.Sum256([]uint8(reportValue)))
		products[report.Object] = reportArtifact
	}

	linkMb.Signatures = []in_toto.Signature{}
	linkMb.Signed = in_toto.Link{
		Type:       "link",
		Name:       "kubesec",
		Materials:  materials,
		Products:   products,
		ByProducts: map[string]interface{}{},
		// FIXME: the command should include whether this is called through the
		// server or a standalone tool.
		Command:     []string{},
		Environment: map[string]interface{}{},
	}

	return linkMb
}

func appendUniqueRule(uniqueRules []RuleRef, newRule RuleRef) []RuleRef {
	if !containsRule(uniqueRules[:], newRule) {
		uniqueRules = append(uniqueRules, newRule)
	}
	return uniqueRules
}

func containsRule(rules []RuleRef, newRule RuleRef) bool {
	for _, rule := range rules {
		if reflect.DeepEqual(rule, newRule) {
			return true
		}
	}
	return false
}

func (rs *Ruleset) generateReport(fileName string, json []byte, schemaDir string) Report {
	report := Report{
		Object:   "Unknown",
		FileName: fileName,
		Score:    0,
		Rules:    make([]RuleRef, 0),
		Scoring: RuleScoring{
			Advise:   make([]RuleRef, 0),
			Passed:   make([]RuleRef, 0),
			Critical: make([]RuleRef, 0),
		},
	}

	report.Object = getObjectName(json)

	// validate resource with kubeval
	cfg := kubeval.NewDefaultConfig()
	cfg.FileName = fileName
	cfg.Strict = true

	if schemaDir != "" {
		cfg.SchemaLocation = "file://" + schemaDir
	} else if _, err := os.Stat("/schemas/kubernetes-json-schema/master/master-standalone"); !os.IsNotExist(err) {
		cfg.SchemaLocation = "file:///schemas"
	}

	results, err := kubeval.Validate(json, cfg)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			report.Message = "This resource is invalid, unknown schema"
		} else {
			report.Message = err.Error()
		}
		return report
	}

	for _, result := range results {
		if len(result.Errors) > 0 {
			for _, desc := range result.Errors {
				report.Message += desc.String() + " "
			}
		} else if result.Kind == "" {
			report.Message += "This resource is invalid, Kubernetes kind not found"
		}
	}

	if len(report.Message) > 0 {
		return report
	}
	report.Valid = true

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

		report.Rules = appendUniqueRule(report.Rules, ruleRef)

		if ruleRef.Containers > 0 {
			if ruleRef.Points >= 0 {
				rs.logger.Debugf("positive score rule matched %v (%v points)", ruleRef.Selector, ruleRef.Points)
				report.Score += ruleRef.Points
				report.Scoring.Passed = append(report.Scoring.Passed, ruleRef)
			}

			if ruleRef.Points < 0 {
				rs.logger.Debugf("negative score rule matched %v (%v points)", ruleRef.Selector, ruleRef.Points)
				report.Score += ruleRef.Points
				report.Scoring.Critical = append(report.Scoring.Critical, ruleRef)
			}
		} else if ruleRef.Points >= 0 {
			rs.logger.Debugf("positive score rule failed %v (%v points)", ruleRef.Selector, ruleRef.Points)
			report.Scoring.Advise = append(report.Scoring.Advise, ruleRef)
		}
	}

	if appliedRules < 1 {
		report.Message = "This resource kind is not supported by kubesec"
	} else if report.Score >= 0 {
		report.Message = fmt.Sprintf("Passed with a score of %v points", report.Score)
	} else {
		report.Message = fmt.Sprintf("Failed with a score of %v points", report.Score)
	}

	// sort results into priority order
	sort.Sort(RuleRefCustomOrder(report.Scoring.Critical))
	sort.Sort(RuleRefCustomOrder(report.Scoring.Passed))
	sort.Sort(RuleRefCustomOrder(report.Scoring.Advise))

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
		ID:         rule.ID,
		Points:     rule.Points,
		Reason:     rule.Reason,
		Selector:   rule.Selector,
		Weight:     rule.Weight,
		Link:       rule.Link,
	}

	ch <- result
}

// getObjectName returns <kind>/<name>.<namespace>
func getObjectName(json []byte) string {
	jq := gojsonq.New().Reader(bytes.NewReader(json))
	if len(jq.Errors()) > 0 {
		return "Unknown"
	}

	kind := jq.Copy().From("kind").Get()
	if kind == nil {
		return "Unknown"
	}
	object := fmt.Sprintf("%v", kind)

	name := jq.Copy().From("metadata.name").Get()
	if name == nil {
		object += "/undefined"
	} else {
		object += fmt.Sprintf("/%v", name)
	}

	namespace := jq.Copy().From("metadata.namespace").Get()
	if namespace == nil {
		object += ".default"
	} else {
		object += fmt.Sprintf(".%v", namespace)
	}

	return object
}

func detectLineBreak(haystack []byte) string {
	windowsLineEnding := bytes.Contains(haystack, []byte("\r\n"))
	if windowsLineEnding && runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
