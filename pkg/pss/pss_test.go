package pss

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const manifestBasePath = "../../test/asset/pss"

func TestRun(t *testing.T) {
	var tests = []struct {
		name                    string // Name of the test
		valid                   bool   // Should the resource to be allowed
		unsupported             bool   // Define if the resource kind is supported
		manifestFilepath        string // Path to manifest file
		profile, profileVersion string // Profile name and version
		expectedReports         int    // Number or reports to expect
	}{
		{
			name:             "Pod satisfies the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pod-restricted.yaml",
			expectedReports:  1,
		},
		{
			name:             "Pod (JSON) satisfies the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pod-restricted.json",
			expectedReports:  1,
		},
		{
			name:             "Pods satisfy the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pods-restricted.yaml",
			expectedReports:  3,
		},
		{
			name:             "Pods (JSON) satisfy the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pods-restricted.json",
			expectedReports:  3,
		},
		{
			name:             "Pod does not satisfy the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            false,
			manifestFilepath: "pod-baseline.yaml",
			expectedReports:  1,
		},
		{
			name:             "Pod satisfies the baseline profile",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pod-baseline.yaml",
			expectedReports:  1,
		},
		{
			name:             "Pod satisfies the privileged profile",
			profile:          "privileged",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "pod-privileged.yaml",
			expectedReports:  1,
		},
		{
			name:             "CronJob satisfies the restricted profile",
			profile:          "restricted",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "cronjob-restricted.yaml",
			expectedReports:  1,
		},
		{
			name:             "Job satisfies the baseline profile",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "job-baseline.yaml",
			expectedReports:  1,
		},
		{
			name:             "DaemonSet does not satisfy the baseline profile",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            false,
			manifestFilepath: "daemonset-privileged.yaml",
			expectedReports:  1,
		},
		{
			name:             "DaemonSet satisfies the privileged profile",
			profile:          "privileged",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "daemonset-privileged.yaml",
			expectedReports:  1,
		},
		{
			name:             "Deployment satisfies the baseline profile",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "deployment-baseline.yaml",
			expectedReports:  1,
		},
		{
			name:             "Deployment from v1beta1 apiVersion is not supported",
			profile:          "privileged", // ignored during test
			profileVersion:   "latest",     // ignored during test
			unsupported:      true,
			manifestFilepath: "deployment-v1beta1.yaml",
			expectedReports:  0,
		},
		{
			name:             "StatefulSet satisfies the baseline profile",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "statefulset-baseline.yaml",
			expectedReports:  1,
		},
		{
			name:             "ServiceAccount is not supported",
			profile:          "restricted", // ignored during test
			profileVersion:   "latest",     // ignored during test
			unsupported:      true,
			manifestFilepath: "serviceaccount.yaml",
			expectedReports:  0,
		},
		{
			name:             "CRDs are not supported",
			profile:          "restricted", // ignored during test
			profileVersion:   "latest",     // ignored during test
			unsupported:      true,
			manifestFilepath: "crds.yaml",
			expectedReports:  0,
		},
		{
			name:             "Only supported resources are scanned and satisfy the baseline policy",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "mixed-baseline.yaml",
			expectedReports:  12,
		},
		{
			name:             "Only supported resources from a full helm template are scanned and satisfy the baseline policy",
			profile:          "baseline",
			profileVersion:   "latest",
			valid:            true,
			manifestFilepath: "helm-vault-baseline.yaml",
			expectedReports:  3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			evaluator, err := NewEvaluator(logger)
			if err != nil {
				t.Fatal(err)
			}

			manifest := readFixtureFile(t, tt.manifestFilepath)
			reports, err := evaluator.Run(tt.manifestFilepath, manifest, tt.profile, tt.profileVersion)
			if err != nil && tt.valid {
				t.Fatal(err)
			}

			// Make sure we get the expected error when the resources
			// do not satisfy the required profile. Unsupported do not return
			// any error as they are simply ignored.
			e := &ProfileNotSatisfiedError{}
			if !tt.valid && !errors.As(err, &e) && !tt.unsupported {
				t.Fatalf("Scanning should have failed with ProfileNotSatisfiedError, got: %v", err)
			}

			if len(reports) != tt.expectedReports {
				t.Errorf("Number of reports invalid, got: %d, want: %d", len(reports), tt.expectedReports)
			}

			for _, report := range reports {
				if report.Valid != tt.valid {
					t.Errorf("Resource validation invalid, got: valid=%t, want: valid=%t", report.Valid, tt.valid)
				}

				if report.FileName != tt.manifestFilepath {
					t.Errorf("Resource filename invalid, got: %s, want: %s", report.FileName, tt.manifestFilepath)

				}
			}
		})
	}

	var tests2 = []struct {
		name                    string // Name of the test
		manifestFilepath        string // Path to manifest file
		profile, profileVersion string // Profile name and version
		expectedReports         int    // Number or reports to expect
	}{
		{
			name:             "Invalid Profile",
			profile:          "notsurethisprofileexists",
			profileVersion:   "latest",
			manifestFilepath: "pod-restricted.yaml",
		},
		{
			name:             "Invalid Profile",
			profile:          "baseline",
			profileVersion:   "notsurethisversionexists",
			manifestFilepath: "pod-baseline.yaml",
		},
	}

	for _, tt := range tests2 {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			evaluator, err := NewEvaluator(logger)
			if err != nil {
				t.Fatal(err)
			}

			manifest := readFixtureFile(t, tt.manifestFilepath)
			reports, err := evaluator.Run(tt.manifestFilepath, manifest, tt.profile, tt.profileVersion)
			if err == nil {
				t.Fatalf("Scan should have failed")
			}

			if len(reports) > tt.expectedReports {
				t.Errorf("Number of reports invalid, got: %d, want: %d", len(reports), tt.expectedReports)
			}
		})
	}

	var tests3 = []struct {
		name                                     string // Name of the test
		manifestFilepath                         string // Path to manifest file
		profile, profileVersion                  string // Profile name and version
		expectedReports                          int    // Number or reports to expect
		invalidObjectReports, validObjectReports int
	}{
		{
			name:                 "Invalid Profile",
			profile:              "baseline",
			profileVersion:       "latest",
			manifestFilepath:     "pods-mixed-profiles.yaml",
			expectedReports:      3,
			invalidObjectReports: 1,
			validObjectReports:   2,
		},
	}

	for _, tt := range tests3 {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			evaluator, err := NewEvaluator(logger)
			if err != nil {
				t.Fatal(err)
			}

			manifest := readFixtureFile(t, tt.manifestFilepath)
			reports, err := evaluator.Run(tt.manifestFilepath, manifest, tt.profile, tt.profileVersion)

			// Make sure we get the expected error when at least one resource
			// does not satisfy the required profile. Unsupported do not return
			// any error as they are simply ignored.
			e := &ProfileNotSatisfiedError{}
			if !errors.As(err, &e) {
				t.Fatalf("Scanning should have failed with ProfileNotSatisfiedError, got: %v", err)
			}

			if len(reports) > tt.expectedReports {
				t.Errorf("Number of reports invalid, got: %d, want: %d", len(reports), tt.expectedReports)
			}

			var valid, invalid int
			for _, report := range reports {
				if report.Valid {
					valid += 1
					continue
				}
				invalid += 1
			}

			if valid != tt.validObjectReports {
				t.Errorf("Number of valid report wrong, got: %d, want: %d", valid, tt.validObjectReports)
			}
			if invalid != tt.invalidObjectReports {
				t.Errorf("Number of invalid report wrong, got: %d, want: %d", invalid, tt.invalidObjectReports)
			}
		})
	}
}

func readFixtureFile(t *testing.T, manifestFilepath string) []byte {
	manifestPath := filepath.Join(manifestBasePath, manifestFilepath)
	manifest, err := os.ReadFile(filepath.Clean(manifestPath))
	if err != nil {
		t.Fatalf("error opening fixture file: %v", err)
	}
	return manifest
}

func TestGetObjectName(t *testing.T) {
	var tests = []struct {
		obj          runtime.Object
		objMeta      metav1.ObjectMeta
		expectedName string
	}{
		{
			expectedName: "Pod/my-name.my-namespace",
			obj: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind: "Pod",
				},
			},
			objMeta: metav1.ObjectMeta{
				Name:      "my-name",
				Namespace: "my-namespace",
			},
		},
		{
			expectedName: "Pod/undefined.default",
			obj: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind: "Pod",
				},
			},
			objMeta: metav1.ObjectMeta{},
		},
	}

	for _, tt := range tests {
		tt := tt

		result := getObjectName(tt.obj, tt.objMeta)
		if result != tt.expectedName {
			t.Errorf("Object name is invalid, got: %s, want: %s", result, tt.expectedName)
		}
	}
}
