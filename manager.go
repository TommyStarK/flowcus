package flowcus

import (
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
	. "github.com/TommyStarK/flowcus/internal/reflect"
)

func newBoxSingleChanTestsManager() *boxSingleChanTestsManager {
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

func (b *boxSingleChanTestsManager) SetTasks(tasks ...tBoxSCTF) {
	for _, task := range tasks {
		b.Set(GetEntityName(task), task)
	}
}

func (b *boxSingleChanTestsManager) StartWorkers(input *Input) {
	var bunch []*Test

	for _, key := range b.Keys() {
		b.Add(1)

		test := newTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				bunch = append(bunch, test)
				wg.Done()
			}()

			task := b.Get(key)
			test.start = time.Now()
			test.caller = GetEntityName(task)
			task.(tBoxSCTF)(test, *input)
			test.finished = true
		}(key, b.WaitGroup, test)
	}

	b.Wait()
	b.Push(&boxSingleChanTestCase{Input: *input, Results: bunch})
}

func newBoxDualChanTestsManager() *boxDualChanTestsManager {
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

func (b *boxDualChanTestsManager) SetTasks(tasks ...tBoxDCTF) {
	for _, task := range tasks {
		b.Set(GetEntityName(task), task)
	}
}

func (b *boxDualChanTestsManager) StartWorkers(input *Input, output *Output) {
	var bunch []*Test

	for _, key := range b.Keys() {
		b.Add(1)

		test := newTest()
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer func() {
				test.duration = time.Since(test.start)
				bunch = append(bunch, test)
				wg.Done()
			}()

			task := b.Get(key)
			test.start = time.Now()
			test.caller = GetEntityName(task)
			task.(tBoxDCTF)(test, *input, *output)
			test.finished = true
		}(key, b.WaitGroup, test)
	}

	b.Wait()
	b.Push(&boxDualChanTestCase{Input: *input, Output: *output, Results: bunch})
}
