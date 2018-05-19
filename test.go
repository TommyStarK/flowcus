package flowcus

import (
	"fmt"
	"runtime"
	"time"
)

type Test struct {
	// *sync.RWMutex
	Caller   string
	Duration time.Duration
	Errors   []string
	Failed   bool
	Finished bool
	Logs     []string
	Start    time.Time
}

func (t *Test) error(s string) {
	// t.Lock()
	// defer t.Unlock()
	t.Errors = append(t.Errors, s)
}

func (t *Test) log(s string) {
	// t.Lock()
	// defer t.Unlock()
	t.Logs = append(t.Logs, s)
}

func (t *Test) Error(args ...interface{}) {
	t.Fail()
	t.error(fmt.Sprintln(args...))
}

func (t *Test) ErrorF(format string, args ...interface{}) {
	t.Fail()
	t.error(fmt.Sprintf(format, args...))
}

func (t *Test) Log(args ...interface{}) {
	t.log(fmt.Sprintln(args...))
}

func (t *Test) LogF(format string, args ...interface{}) {
	t.log(fmt.Sprintf(format, args...))
}

func (t *Test) Fail() {
	// t.Lock()
	// defer t.Unlock()
	t.Failed = true
}

func (t *Test) FailNow() {
	t.Fail()
	t.Finished = true
	runtime.Goexit()
}

func (t *Test) Fatal(args ...interface{}) {
	t.Error(args...)
	t.FailNow()
}

func (t *Test) FatalF(format string, args ...interface{}) {
	t.ErrorF(format, args...)
	t.FailNow()
}
