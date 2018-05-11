package flowcus

import (
	"fmt"
	"time"
)

type Flow struct {
	Data     interface{}
	duration time.Duration
	sample   interface{}
	success  bool
	tester   string
}

func (f *Flow) Empty() bool {
	return *f == (Flow{})
}

func (f *Flow) String() string {
	return fmt.Sprintf("%p - %#v", f, f)
}
