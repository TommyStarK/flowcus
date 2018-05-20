package flowcus

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type TestExported struct {
	Caller   string
	Start    time.Time
	Duration time.Duration
	Finished bool
	Skipped  bool
	Success  bool
	Errors   []string
	Logs     []string
}

func NewTest() *Test {
	return &Test{
		"",
		time.Time{},
		0,
		false,
		false,
		false,
		false,
		[]string{},
		[]string{},
		&sync.RWMutex{},
	}
}

type Test struct {
	caller   string
	start    time.Time
	duration time.Duration
	failed   bool
	finished bool
	skipped  bool
	success  bool
	errors   []string
	logs     []string
	*sync.RWMutex
}

func (t *Test) error(s string) {
	t.Lock()
	defer t.Unlock()
	t.errors = append(t.errors, s)
}

func (t *Test) log(s string) {
	t.Lock()
	defer t.Unlock()
	t.logs = append(t.logs, s)
}

func (t *Test) skip() {
	t.Lock()
	defer t.Unlock()
	t.skipped = true
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
	t.finished = true
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

func (t *Test) Skip(args ...interface{}) {
	t.Log(args...)
	t.skip()
}

func (t *Test) SkipF(format string, args ...interface{}) {
	t.LogF(format, args...)
	t.skip()
}

func (t *Test) SkipNow() {
	t.skip()
	t.finished = true
	runtime.Goexit()
}

func (t *Test) Skipped() bool {
	t.RLock()
	skipped := t.skipped
	t.RUnlock()
	return skipped
}
