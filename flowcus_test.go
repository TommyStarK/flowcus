package flowcus

import (
	"testing"
)

func TestNewFlowcus(t *testing.T) {
	flowcus := NewFlowcus()

	if flowcus.producer != nil {
		t.Errorf("producer should be nil")
	}

	if flowcus.consumer != nil {
		t.Errorf("consumer should be nil")
	}

	if flowcus.waitGrp == nil {
		t.Errorf("waitGrp should be of type *sync.WaitGroup")
	}

	if flowcus.tests.Len() > 0 {
		t.Errorf("no test should be stored")
	}

	if flowcus.jobs.Len() > 0 {
		t.Errorf("no job should be waiting ")
	}

	if flowcus.event == nil {
		t.Errorf("event should be of type chan *flowcus.Event")
	}

	if flowcus.revent == nil {
		t.Errorf("revent should be of type chan *flowcus.Revent")
	}
}

func TestBasicFlow(t *testing.T) {
	f := NewFlowcus()

	test := func(omap *OrderedMap, data interface{}) (interface{}, error) {
		return data.(int), nil
	}

	f.Producer(func(com chan<- *Event) {
		defer func() {
			close(com)
		}()

		for index := 0; index < 50; index++ {
			com <- &Event{
				Id: index,
			}
		}
	})

	f.Consumer(func(com chan<- *Revent) {
		defer func() {
			close(com)
		}()

		for index := 0; index < 50; index++ {
			com <- &Revent{
				Data: index,
				Test: test,
			}
		}
	})

	f.Start()

	report := f.Report()
	if report == nil {
		t.Errorf("report should not be nil")
	}

	if len(report.Tests) < 50 {
		t.Errorf("report.Tests should contain 50 tests")
	}

	if report.Errors != nil {
		t.Errorf("report.Errors should be nil")
	}
}
