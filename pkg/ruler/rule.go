package ruler

import (
	"bytes"
	"fmt"

	"github.com/thedevsaddam/gojsonq/v2"
)

type NotSupportedError struct {
	Kind string
}

func (e *NotSupportedError) Error() string {
	return fmt.Sprintf("rule does not apply to kind %s", e.Kind)
}

type Rule struct {
	ID        string           `json:"id" yaml:"id"`
	Selector  string           `json:"selector" yaml:"selector"`
	Reason    string           `json:"reason" yaml:"reason"`
	Link      string           `json:"link,omitempty" yaml:"link,omitempty"`
	Kinds     []string         `json:"kinds" yaml:"kinds"`
	Points    int              `json:"points" yaml:"points"`
	Advise    int              `json:"advise" yaml:"advise"`
	Predicate func([]byte) int `json:"-" yaml:"-"`
}

// Eval executes the predicate if the kind matches the rule
func (r *Rule) Eval(json []byte) (int, error) {
	jq := gojsonq.New().Reader(bytes.NewReader(json)).From("kind")
	if jq.Error() != nil {
		return 0, jq.Error()
	}

	kind := fmt.Sprintf("%s", jq.Get())

	var match bool
	for _, k := range r.Kinds {
		if k == kind {
			match = true
			break
		}
	}

	if match {
		count := r.Predicate(json)
		return count, nil
	} else {
		return 0, &NotSupportedError{Kind: kind}
	}
}
