package jira

type Transition struct {
	Backlog   int
	ReadyToDo int
	Progress  int
	Done      int
	Review    int
	Homol     int
}

var Transitions = Transition{
	Backlog:   11,
	ReadyToDo: 21,
	Progress:  31,
	Done:      41,
	Review:    51,
	Homol:     81,
}
