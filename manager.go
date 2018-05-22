package flowcus

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
)

//
// Exploratory Box Tests Manager
//
func NewExploratoryBoxTestsManager() *exploratoryBoxTestsManager {
	return &exploratoryBoxTestsManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}

type exploratoryBoxTestCase struct {
	Input   Input
	Results []*Test
}

type exploratoryBoxTestsManager struct {
	cases *Fifo
	tests *OrderedMap
	mutex *sync.Mutex
	wg    *sync.WaitGroup
}

func (e *exploratoryBoxTestsManager) SetTasks(tasks ...BoxETF) {
	for _, task := range tasks {
		e.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (e *exploratoryBoxTestsManager) StartWorkers(input *Input) {
	bunch := make([]*Test, 0)

	for _, key := range e.tests.Keys() {
		e.wg.Add(1)

		test := NewTest()
		go func(key interface{}, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				if r := recover(); r != nil {
					switch r.(type) {
					case runtime.Error:
						panic(r)
					default:
						test.Fail()
						test.Error(r)
					}
				}

				e.mutex.Lock()
				bunch = append(bunch, test)
				e.mutex.Unlock()
				e.wg.Done()
			}()

			task := e.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxETF)(test, *input)
			test.finished = true
		}(key, test)
	}

	e.wg.Wait()
	e.cases.Push(&exploratoryBoxTestCase{Input: *input, Results: bunch})
}

//
// Linear Box Tests Manager
//
func NewLinearBoxTestsManager() *linearBoxTestsManager {
	return &linearBoxTestsManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}

type linearBoxTestCase struct {
	Input   Input
	Output  Output
	Results []*Test
}

type linearBoxTestsManager struct {
	cases *Fifo
	tests *OrderedMap
	mutex *sync.Mutex
	wg    *sync.WaitGroup
}

func (l *linearBoxTestsManager) SetTasks(tasks ...BoxLTF) {
	for _, task := range tasks {
		l.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (l *linearBoxTestsManager) StartWorkers(input *Input, output *Output) {
	bunch := make([]*Test, 0)

	for _, key := range l.tests.Keys() {
		l.wg.Add(1)

		test := NewTest()
		go func(key interface{}, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				if r := recover(); r != nil {
					switch r.(type) {
					case runtime.Error:
						panic(r)
					default:
						test.Fail()
						test.Error(r)
					}
				}

				l.mutex.Lock()
				bunch = append(bunch, test)
				l.mutex.Unlock()
				l.wg.Done()
			}()

			task := l.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxLTF)(test, *input, *output)
			test.finished = true
		}(key, test)
	}

	l.wg.Wait()
	l.cases.Push(&linearBoxTestCase{Input: *input, Output: *output, Results: bunch})
}

//
// Non Linear Box Tests Manager
//
func NewNonLinearBoxTestsManager() *nonlinearBoxTestsManager {
	return &nonlinearBoxTestsManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}

type nonlinearBoxTestCase struct {
	Inputs  []Input
	Outputs []Output
	Results []*Test
}

type nonlinearBoxTestsManager struct {
	cases *Fifo
	tests *OrderedMap
	mutex *sync.Mutex
	wg    *sync.WaitGroup
}

func (n *nonlinearBoxTestsManager) SetTasks(tasks ...BoxNLTF) {
	for _, task := range tasks {
		n.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (n *nonlinearBoxTestsManager) StartWorkers(inputs []Input, outputs []Output) {
	bunch := make([]*Test, 0)

	for _, key := range n.tests.Keys() {
		n.wg.Add(1)

		test := NewTest()
		go func(key interface{}, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				if r := recover(); r != nil {
					switch r.(type) {
					case runtime.Error:
						panic(r)
					default:
						test.Fail()
						test.Error(r)
					}
				}

				n.mutex.Lock()
				bunch = append(bunch, test)
				n.mutex.Unlock()
				n.wg.Done()
			}()

			task := n.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxNLTF)(test, inputs, outputs)
			test.finished = true
		}(key, test)
	}

	n.wg.Wait()
	n.cases.Push(&nonlinearBoxTestCase{Inputs: inputs, Outputs: outputs, Results: bunch})
}
