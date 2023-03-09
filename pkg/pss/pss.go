package pss

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/controlplaneio/kubesec/v2/pkg/util"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/pod-security-admission/api"
	policy "k8s.io/pod-security-admission/policy"
)

// Evaluator holds the interfaces to decode manifests and validate a PSS policy.
type Evaluator struct {
	logger *zap.SugaredLogger
	policy.Evaluator
	runtime.Decoder
}

// NewEvaluator returns a new PSS Evaluator.
func NewEvaluator(logger *zap.SugaredLogger) (*Evaluator, error) {
	evaluator, err := policy.NewEvaluator(policy.DefaultChecks())
	if err != nil {
		return nil, fmt.Errorf("Unable to create policy evaluator: %w", err)
	}

	return &Evaluator{
		Evaluator: evaluator,
		logger:    logger,
		Decoder:   scheme.Codecs.UniversalDeserializer(),
	}, nil
}

// Reports holds the info about a PSS check against a single object.
type Report struct {
	Object          string           `json:"object" yaml:"object"`
	Valid           bool             `json:"valid" yaml:"valid"`
	FileName        string           `json:"fileName,omitempty" yaml:"fileName,omitempty"`
	Policy          string           `json:"policy" yaml:"policy"`
	PolicyVersion   string           `json:"policyVersion" yaml:"policyVersion"`
	ForbiddenChecks []ForbiddenCheck `json:"forbiddenChecks" yaml:"forbiddenChecks"`
}

// ForbiddenCheck holds the checks forbidden against a policy.
type ForbiddenCheck struct {
	Reason string `json:"reason" yaml:"reason"`
	Detail string `json:"detail" yaml:"detail"`
}

// PolicyNotSatisfiedError is the occurs happening when a resource does not
// satisfies a policy.
type PolicyNotSatisfiedError struct {
	Policy        string
	PolicyVersion string
}

// Error satisfies the Error interface.
func (e *PolicyNotSatisfiedError) Error() string {
	return fmt.Sprintf("One or more resources do not satisfy the PSS policy: %s/%s", e.Policy, e.PolicyVersion)
}

// Run checks the content of a manifest as bytes (fileName is only used for the report)
// against a versioned PSS policy and returns a Report with the results.
func (e *Evaluator) Run(fileName string, fileBytes []byte, policy, policyVersion string) ([]Report, error) {
	reports := make([]Report, 0)

	lvl, err := api.ParseLevel(policy)
	if err != nil {
		return reports, fmt.Errorf("Invalid policy: %s", policy)
	}
	ver, err := api.ParseVersion(policyVersion)
	if err != nil {
		return reports, fmt.Errorf("Invalid policy version: %s", policyVersion)
	}
	levelVersion := api.LevelVersion{
		Level:   lvl,
		Version: ver,
	}

	objects := util.GetObjectsFromManifest(fileBytes)
	if len(objects) == 0 {
		e.logger.Infof("Unable to find any YAML or JSON resources in file: %s", fileName)
		return reports, nil
	}

	for _, obj := range objects {
		o, err := e.decode(bytes.TrimSpace(obj))
		if err != nil {
			// We don't care about objects that can't be decoded
			continue
		}

		report, err := e.eval(o, levelVersion)
		if err != nil {
			return reports, err
		}
		if reflect.DeepEqual(report, Report{}) {
			continue
		}

		report.FileName = fileName
		reports = append(reports, report)
	}

	for _, report := range reports {
		if !report.Valid {
			return reports, &PolicyNotSatisfiedError{
				Policy:        policy,
				PolicyVersion: policyVersion,
			}
		}
	}

	return reports, nil
}

func (e *Evaluator) eval(obj runtime.Object, lv api.LevelVersion) (Report, error) {
	var (
		objName     string
		podMetadata metav1.ObjectMeta
		podSpec     corev1.PodSpec
	)

	// Make sure the input object is a supported resource
	switch o := obj.(type) {
	case *corev1.Pod:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.ObjectMeta
		podSpec = o.Spec
	case *appsv1.DaemonSet:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.Spec.Template.ObjectMeta
		podSpec = o.Spec.Template.Spec
	case *appsv1.Deployment:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.Spec.Template.ObjectMeta
		podSpec = o.Spec.Template.Spec
	case *appsv1.StatefulSet:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.Spec.Template.ObjectMeta
		podSpec = o.Spec.Template.Spec
	case *batchv1.CronJob:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.Spec.JobTemplate.Spec.Template.ObjectMeta
		podSpec = o.Spec.JobTemplate.Spec.Template.Spec
	case *batchv1.Job:
		objName = getObjectName(o, o.ObjectMeta)
		podMetadata = o.Spec.Template.ObjectMeta
		podSpec = o.Spec.Template.Spec
	default:
		e.logger.Debugf("Resource not supported, validation skipped: %s",
			obj.DeepCopyObject().GetObjectKind().GroupVersionKind().String())
		return Report{}, nil
	}

	// Evaluate
	results := e.EvaluatePod(lv, &podMetadata, &podSpec)
	aggregate := policy.AggregateCheckResults(results)
	report := Report{
		Object:        objName,
		Policy:        string(lv.Level),
		PolicyVersion: lv.Version.String(),
		Valid:         aggregate.Allowed,
	}

	// Report the forbidden checks
	for i := range results {
		if !results[i].Allowed {
			report.ForbiddenChecks = append(report.ForbiddenChecks, ForbiddenCheck{
				Detail: results[i].ForbiddenDetail,
				Reason: results[i].ForbiddenReason,
			})
		}
	}

	return report, nil
}

// decode attempts to deserialize an object as bytes (from a manifest).
func (e *Evaluator) decode(data []byte) (runtime.Object, error) {
	obj, _, err := e.Decode(data, nil, nil)
	// These errors happen for non Kubernetes objects or OpenShift objects.
	// From runtime: notRegisteredError, isMissingVersion, missingKind.
	if err != nil {
		e.logger.Debugf("Unable to decode resource: %s", err.Error())
		return nil, err
	}

	return obj, nil
}

// getObjectName returns a standardized name with the format: <kind>/<name>.<namespace>
func getObjectName(obj runtime.Object, objMeta metav1.ObjectMeta) string {
	if obj == nil {
		return ""
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind

	name := objMeta.Name
	if name == "" {
		name = "undefined"
	}

	namespace := objMeta.Namespace
	if namespace == "" {
		namespace = "default"
	}

	return kind + "/" + name + "." + namespace
}
