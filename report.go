package flowcus

import (
	"time"
)

type Error struct {
	Date string
	Err  string
}

type Test struct {
	Id       interface{}
	Label    string
	Success  bool
	Duration time.Duration
	Tester   string
	Sample   interface{}
}

func (t *Test) Empty() bool {
	return *t == (Test{})
}

type Report struct {
	Date     string
	Version  float64
	Number   int
	Coverage float64
	Duration time.Duration
	Type     string
	Tests    []*Test
	Errors   []*Error
}

func (r *Report) Empty() bool {
	return len(r.Tests) == 0
}
