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

	if flowcus.mutex == nil {
		t.Errorf("mutex should be of type *sync.Mutex")
	}

	if flowcus.wait == nil {
		t.Errorf("wait should be of type *sync.WaitGroup")
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
