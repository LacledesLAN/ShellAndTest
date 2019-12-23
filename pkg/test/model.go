package test

// testCriteria contains all the criteria and information needed from the JSON file.
type testCriteria struct {
	ShouldEcho []struct {
		Command    string `json:"command"`
		ShouldHave string `json:"should-have"`
	} `json:"should-echo"`
	ShouldHave []string `json:"should-have"`
	ShouldLack []string `json:"should-lack"`
	Target     struct {
		Execute         string   `json:"execute"`
		PostTasks       []string `json:"post-tasks"`
		PreTasks        []string `json:"pre-tasks"`
		ShouldEchoDelay int      `json:"should-echo-delay"`
		Timeout         int      `json:"timeout"`
	} `json:"target"`
}
