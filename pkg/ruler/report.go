package ruler

type Report struct {
	Object  string      `json:"object"`
	Valid   bool        `json:"valid"`
	Message string      `json:"message,omitempty"`
	Score   int         `json:"score"`
	Scoring RuleScoring `json:"scoring,omitempty"`
}

type RuleScoring struct {
	Critical []RuleRef `json:"critical,omitempty"`
	Advise   []RuleRef `json:"advise,omitempty"`
}

type RuleRef struct {
	Selector   string `json:"selector"`
	Reason     string `json:"reason"`
	Weight     int    `json:"weight,omitempty"`
	Link       string `json:"href,omitempty"`
	Containers int    `json:"-"`
	Points     int    `json:"-"`
}
