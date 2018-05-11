package flowcus

import (
	"time"
)

type Test struct {
	Id       interface{}
	Duration time.Duration
	Sample   interface{}
	Success  bool
	Tester   string
}

type Report struct {
	Date     string
	Duration time.Duration
	Type     string
	Version  float64
	Tests    []*Test
}
