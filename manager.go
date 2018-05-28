package flowcus

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
)

func acquireTest(test *Test, bunch *[]*Test, mutex *sync.Mutex, waitGroup *sync.WaitGroup) {
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
	waitGroup.Done()
}

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
			defer acquireTest(test, &bunch, e.mutex, e.wg)

			task := e.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxETF)(test, *input)
			test.finished = true
		}(key, test)
	}

	e.wg.Wait()
	e.cases.Push(bunch)
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
			defer acquireTest(test, &bunch, l.mutex, l.wg)

			task := l.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxLTF)(test, *input, *output)
			test.finished = true
		}(key, test)
	}

	l.wg.Wait()
	l.cases.Push(bunch)
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
			defer acquireTest(test, &bunch, n.mutex, n.wg)

			task := n.tests.Get(key)
			test.start = time.Now()
			test.caller = runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
			task.(BoxNLTF)(test, inputs, outputs)
			test.finished = true
		}(key, test)
	}

	n.wg.Wait()
	n.cases.Push(bunch)
}
