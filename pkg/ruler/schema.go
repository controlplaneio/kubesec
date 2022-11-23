package ruler

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yannh/kubeconform/pkg/validator"
)

// SchemaConfig hold the configuration of the schema validaton.
type SchemaConfig struct {
	// DisableValidation disables the validation of the manifests against
	// Kubernetes JSON schema. Set to true when the source manifests
	// comes directly from the cluster (e.g: webhook, kubectl plugin).
	DisableValidation bool

	// Locations defines the locations of the schemas. This follows the
	// same logic as the -schema-location flag from kubeconform.
	Locations []string

	// ValidatorOpts are the options from kubeconform validator.
	ValidatorOpts validator.Opts
}

func NewDefaultSchemaConfig() SchemaConfig {
	return SchemaConfig{
		ValidatorOpts: validator.Opts{
			Strict: true,
		},
	}
}

// validateSchema validates the json schema of the resource
// using kubeconform and updates the provided Report.
func validateSchema(report Report, json []byte, schemaConfig SchemaConfig) Report {
	v, err := validator.New(schemaConfig.Locations, schemaConfig.ValidatorOpts)
	if err != nil {
		report.Message += fmt.Sprintf("failed initializing validator: %s", err)
		return report
	}

	f := io.NopCloser(bytes.NewReader(json))
	for _, res := range v.Validate(report.FileName, f) {
		// A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status == validator.Invalid {
			report.Message += res.Err.Error() + "\n"
		}
		if res.Status == validator.Error {
			report.Message += res.Err.Error()
		}
	}

	return report
}
