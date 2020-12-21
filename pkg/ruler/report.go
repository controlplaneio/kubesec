package ruler

type Reports []Report

type Report struct {
	Object   string      `json:"object"`
	Valid    bool        `json:"valid"`
	FileName string      `json:"fileName"`
	Rules    []RuleRef   `json:"-"`
	Message  string      `json:"message,omitempty"`
	Score    int         `json:"score"`
	Scoring  RuleScoring `json:"scoring,omitempty"`
}

type RuleScoring struct {
	Critical []RuleRef `json:"critical,omitempty"`
	Passed   []RuleRef `json:"passed,omitempty"`
	Advise   []RuleRef `json:"advise,omitempty"`
}

type RuleRef struct {
	ID         string `json:"id"`
	Selector   string `json:"selector"`
	Reason     string `json:"reason"`
	Weight     int    `json:"weight,omitempty"`
	Link       string `json:"href,omitempty"`
	Containers int    `json:"-"`
	Points     int    `json:"points"`
}

// This implements a custom sort interface (Len, Swap, Less) for the report listing.
// Each scan can produce a different ordering of the reported tests. To have a single
// deterministic report response for the same input requires sort to never draw.
// Assumption below is: the combination of points, then selector text, should be unique
// This is applied to the output of scan for each of the Critical and Advisory lists.

type RuleRefCustomOrder []RuleRef

func (rr RuleRefCustomOrder) Len() int { return len(rr) }

func (rr RuleRefCustomOrder) Swap(i, j int) { rr[i], rr[j] = rr[j], rr[i] }

func (rr RuleRefCustomOrder) Less(i, j int) bool {
	if rr[i].Points != rr[j].Points {
		// no integer absolute fn in golang
		if rr[i].Points > 0 || rr[j].Points > 0 {
			return rr[i].Points > rr[j].Points
		}
		return rr[i].Points < rr[j].Points
	}
	return rr[i].Selector < rr[j].Selector
}
