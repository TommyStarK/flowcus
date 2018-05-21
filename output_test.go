package flowcus

import "testing"

func TestOutput(t *testing.T) {
	output := &Output{}

	if !output.Empty() {
		t.Error("output struct should be empty")
		t.Log(output.String())
	}

	output.Data = []byte("data")

	if output.Empty() {
		t.Error("output struct should not be empty")
		t.Log(output.String())
	}
}
