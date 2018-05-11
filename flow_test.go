package flowcus

import (
	"testing"
)

func TestFlow(t *testing.T) {
	flow := &Flow{}

	if !flow.Empty() {
		t.Errorf("flow should be empty")
	}

	flow.Data = "test"

	if flow.Empty() {
		t.Errorf("flow should not be empty")
	}

	str := flow.String()
	if len(str) == 0 {
		t.Errorf("calling flow.String() should return a non empty string")
	} else {
		t.Log(str)
	}
}
