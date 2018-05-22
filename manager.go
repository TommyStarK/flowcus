package flowcus

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
)

func acquireTest(wg *sync.WaitGroup, test *Test, bunch *[]*Test) {
	mutex := &sync.Mutex{}

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

	mutex.Lock()
	*bunch = append(*bunch, test)
	mutex.Unlock()
	wg.Done()
}

//
// Exploratory Box Tests Manager
//
func NewExploratoryBoxTestsManager() *exploratoryBoxTestsManager {
	return &exploratoryBoxTestsManager{
		NewFifo(),
		NewOrderedMap(),
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
	wg    *sync.WaitGroup
}

func (e *exploratoryBoxTestsManager) SetTasks(tasks ...BoxETF) {
	for _, task := range tasks {
		e.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (e *exploratoryBoxTestsManager) StartWorkers(input *Input) {
	// var bunch []*Test
	bunch := make([]*Test, 0)

	for _, key := range e.tests.Keys() {
		e.wg.Add(1)

		test := NewTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			// defer acquireTest(wg, test, &bunch)
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
				bunch = append(bunch, test)
				wg.Done()
			}()

			task := e.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxETF)(test, *input)
			test.finished = true
		}(key, e.wg, test)
		<-time.After(100 * time.Microsecond)
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
	wg    *sync.WaitGroup
}

func (l *linearBoxTestsManager) SetTasks(tasks ...BoxLTF) {
	for _, task := range tasks {
		l.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (l *linearBoxTestsManager) StartWorkers(input *Input, output *Output) {
	var bunch []*Test

	for _, key := range l.tests.Keys() {
		l.wg.Add(1)

		test := NewTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer acquireTest(wg, test, &bunch)

			task := l.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxLTF)(test, *input, *output)
			test.finished = true
		}(key, l.wg, test)
		<-time.After(100 * time.Microsecond)
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
	wg    *sync.WaitGroup
}

func (n *nonlinearBoxTestsManager) SetTasks(tasks ...BoxNLTF) {
	for _, task := range tasks {
		n.tests.Set(runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), task)
	}
}

func (n *nonlinearBoxTestsManager) StartWorkers(inputs []Input, outputs []Output) {
	var bunch []*Test

	for _, key := range n.tests.Keys() {
		n.wg.Add(1)

		test := NewTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer acquireTest(wg, test, &bunch)

			task := n.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxNLTF)(test, inputs, outputs)
			test.finished = true
		}(key, n.wg, test)
		<-time.After(100 * time.Microsecond)
	}

	n.wg.Wait()
	n.cases.Push(&nonlinearBoxTestCase{Inputs: inputs, Outputs: outputs, Results: bunch})
}
