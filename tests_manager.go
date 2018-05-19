package flowcus

import (
	"sync"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
	. "github.com/TommyStarK/flowcus/internal/ordered_map"
	. "github.com/TommyStarK/flowcus/internal/reflect"
)

func _bboxManager() *bboxManager {
	return &bboxManager{
		NewFifo(),
		NewOrderedMap(),
		&sync.WaitGroup{},
	}
}

type bboxManager struct {
	*Fifo
	*OrderedMap
	*sync.WaitGroup
}

type bboxTestCase struct {
	Input   *Input
	Output  *Output
	Results []*Test
}

func (b *bboxManager) SetTasks(tasks ...tFunc) {
	for _, task := range tasks {
		b.Set(GetEntityName(task), task)
	}
}

func (b *bboxManager) StartWorkers(input *Input, output *Output) {
	var bunch []*Test

	for _, key := range b.Keys() {
		b.Add(1)

		var test Test
		go func(key interface{}, wg *sync.WaitGroup, test *Test) {
			defer func() {
				test.Duration = time.Since(test.Start)

				wg.Done()
			}()

			task := b.Get(key)
			test.Start = time.Now()
			test.Caller = GetEntityName(task)
			task.(tFunc)(test, input, output)
			test.Finished = true
			bunch = append(bunch, test)
		}(key, b.WaitGroup, &test)
	}

	b.Wait()
	b.Push(&bboxTestCase{Input: input, Output: output, Results: bunch})
}
