package flowcus

import "testing"

func TestEvent(t *testing.T) {
	event := Event{}

	if !event.Empty() {
		t.Errorf("event should be empty")
	}

	event.Id = 1

	if event.Empty() {
		t.Errorf("event should not be empty")
	}

	str := event.String()
	if len(str) == 0 {
		t.Errorf("calling event.String() should return a non empty string")
	} else {
		t.Log(str)
	}
}
