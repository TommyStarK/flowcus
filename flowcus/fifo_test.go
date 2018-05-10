package flowcus

import "testing"

func TestNewFifo(t *testing.T) {
	fifo := NewFifo()

	if fifo.head != nil {
		t.Errorf("head should ne bil")
	}

	if fifo.tail != nil {
		t.Errorf("tail should be nil")
	}

	if fifo.mutex == nil {
		t.Errorf("mutex should be of type *sync.Mutex")
	}

	if fifo.len != 0 {
		t.Errorf("len should be equal to 0")
	}
}

func TestFifoUseCases(t *testing.T) {
	fifo := NewFifo()

	if shouldBeNil := fifo.Pop(); shouldBeNil != nil {
		t.Errorf("fifo is empty, calling Pop() should return nil")
	}

	fifo.Push("test one")

	if len := fifo.Len(); len != 1 {
		t.Errorf("after pushing the first element, len should be equal to 1")
	}

	fifo.Push("test two")
	fifo.Push("test three")

	if len := fifo.Len(); len != 3 {
		t.Errorf("len should be equal to 3")
	}

	elemOne := fifo.Pop()
	if elemOne.(string) != "test one" {
		t.Errorf("first element popped should be equal to 'test one'")
	}

	elemTwo := fifo.Pop()
	if elemTwo.(string) != "test two" {
		t.Errorf("second element popped should be equal to 'test two'")
	}

	elemThree := fifo.Pop()
	if elemThree.(string) != "test three" {
		t.Errorf("third element popped should be equal to 'test three'")
	}
}
