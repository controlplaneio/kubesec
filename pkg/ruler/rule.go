package ruler

import (
	"bytes"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
)

type Rule struct {
	Selector  string
	Title     string
	Reason    string
	Link      string
	Kinds     []string
	Points    int
	Weight    int
	Advise    int
	Predicate func([]byte) int
}

// Eval executes the predicate if the kind matches the rule
func (r *Rule) Eval(json []byte) (bool, error) {
	jq := gojsonq.New().Reader(bytes.NewReader(json)).From("kind")
	if jq.Error() != nil {
		return true, jq.Error()
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
		return count < 1, nil
	} else {
		return true, fmt.Errorf("rule does not apply to kind %s", kind)
	}
}
