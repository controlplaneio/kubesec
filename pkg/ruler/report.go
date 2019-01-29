package ruler

type Report struct {
	Error   string      `json:"error,omitempty"`
	Score   int         `json:"score"`
	Scoring RuleScoring `json:"scoring"`
}

type RuleScoring struct {
	Critical []RuleRef `json:"critical"`
	Advise   []RuleRef `json:"advise"`
}

type RuleRef struct {
	Selector string `json:"selector"`
	Reason   string `json:"reason"`
	Weight   int    `json:"weight"`
	Link     string `json:"href,omitempty"`
}
