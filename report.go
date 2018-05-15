package flowcus

import (
	"time"
)

type Test struct {
	Id       interface{}
	Duration time.Duration
	Label    string
	Sample   interface{}
	Success  bool
	Tester   string
}

func (t *Test) Empty() bool {
	return *t == (Test{})
}

type Report struct {
	Coverage float64
	Date     string
	Duration time.Duration
	Errors   []string
	Number   int
	Type     string
	Version  float64
	Tests    []*Test
}

func (r *Report) Empty() bool {
	return len(r.Tests) == 0
}
