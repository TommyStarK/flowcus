package flowcus

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
)

func NewBoxSingleChanTestsManager() *boxSingleChanTestsManager {
	return &boxSingleChanTestsManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.WaitGroup{},
	}
}

type boxSingleChanTestCase struct {
	Input   Input
	Results []*Test
}

type boxSingleChanTestsManager struct {
	*Fifo
	*OrderedMap
	*sync.WaitGroup
}

func (b *boxSingleChanTestsManager) SetTasks(tasks ...BoxSCTF) {
	for _, task := range tasks {
		b.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (b *boxSingleChanTestsManager) StartWorkers(input *Input) {
	var bunch []*Test

	for _, key := range b.Keys() {
		b.Add(1)

		test := NewTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				bunch = append(bunch, test)
				wg.Done()
			}()

			task := b.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxSCTF)(test, *input)
			test.finished = true
		}(key, b.WaitGroup, test)
		<-time.After(100 * time.Microsecond)
	}

	b.Wait()
	b.Push(&boxSingleChanTestCase{Input: *input, Results: bunch})
}

func NewBoxDualChanTestsManager() *boxDualChanTestsManager {
	return &boxDualChanTestsManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.WaitGroup{},
	}
}

type boxDualChanTestCase struct {
	Input   Input
	Output  Output
	Results []*Test
}

type boxDualChanTestsManager struct {
	*Fifo
	*OrderedMap
	*sync.WaitGroup
}

func (b *boxDualChanTestsManager) SetTasks(tasks ...BoxDCTF) {
	for _, task := range tasks {
		b.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (b *boxDualChanTestsManager) StartWorkers(input *Input, output *Output) {
	var bunch []*Test

	for _, key := range b.Keys() {
		b.Add(1)

		test := NewTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				bunch = append(bunch, test)
				wg.Done()
			}()

			task := b.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxDCTF)(test, *input, *output)
			test.finished = true
		}(key, b.WaitGroup, test)
		<-time.After(100 * time.Microsecond)
	}

	b.Wait()
	b.Push(&boxDualChanTestCase{Input: *input, Output: *output, Results: bunch})
}
