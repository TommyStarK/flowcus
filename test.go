package flowcus

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func newTest() *Test {
	return &Test{
		"",
		time.Time{},
		0,
		false,
		false,
		false,
		[]string{},
		[]string{},
		&sync.RWMutex{},
	}
}

type Test struct {
	Caller   string
	Start    time.Time
	Duration time.Duration
	failed   bool
	Finished bool
	Success  bool
	Errors   []string
	Logs     []string
	*sync.RWMutex
}

func (t *Test) error(s string) {
	t.Lock()
	defer t.Unlock()
	t.Errors = append(t.Errors, s)
}

func (t *Test) log(s string) {
	t.Lock()
	defer t.Unlock()
	t.Logs = append(t.Logs, s)
}

func (t *Test) Error(args ...interface{}) {
	t.error(fmt.Sprintln(args...))
	t.Fail()
}

func (t *Test) ErrorF(format string, args ...interface{}) {
	t.error(fmt.Sprintf(format, args...))
	t.Fail()
}

func (t *Test) Log(args ...interface{}) {
	t.log(fmt.Sprintln(args...))
}

func (t *Test) LogF(format string, args ...interface{}) {
	t.log(fmt.Sprintf(format, args...))
}

func (t *Test) Fail() {
	if t.Failed() {
		return
	}

	t.Lock()
	defer t.Unlock()
	t.failed = true
}
func (t *Test) Failed() bool {
	t.RLock()
	failed := t.failed
	t.RUnlock()
	return failed
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
