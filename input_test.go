package flowcus

import "testing"

func TestInput(t *testing.T) {
	input := &Input{}

	if !input.Empty() {
		t.Error("input struct should be empty")
		t.Log(input.String())
	}

	input.Data = []int{1, 2, 3}
	input.Id = "this is an id"
	input.Expected = 1
	input.Label = "should fail"

	if input.Empty() {
		t.Error("input struct should not be empty")
		t.Log(input.String())
	}
}
