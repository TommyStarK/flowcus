package flowcus

import (
	"sync"
)

func NewFifo() *Fifo {
	return &Fifo{
		0,
		nil,
		nil,
		&sync.Mutex{},
	}
}

type element struct {
	data interface{}
	next *element
}

type Fifo struct {
	len   int
	head  *element
	tail  *element
	mutex *sync.Mutex
}

func (f *Fifo) Len() int {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.len
}

func (f *Fifo) Push(data interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	elem := &element{data: data, next: nil}
	if f.tail == nil {
		f.head = elem
		f.tail = elem
	} else {
		last := f.tail
		last.next = elem
		f.tail = elem
	}
	f.len++
}

func (f *Fifo) Pop() interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.head == nil {
		return nil
	}

	head := f.head
	if head.next == nil {
		f.head = nil
		f.tail = nil
	} else {
		f.head = head.next
	}

	f.len--
	return head.data
}
