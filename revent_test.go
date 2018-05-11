package flowcus

import "testing"

func TestRevent(t *testing.T) {
	revent := &Revent{}

	if !revent.Empty() {
		t.Errorf("revent should be empty")
	}

	revent.Data = "test"

	if revent.Empty() {
		t.Errorf("revent should not be empty")
	}

	str := revent.String()
	if len(str) == 0 {
		t.Errorf("calling revent.String() should return a non empty string")
	} else {
		t.Log(str)
	}
}
