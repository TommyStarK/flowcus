package flowcus

import (
	"fmt"
)

type Event struct {
	Id interface{}
}

func (e *Event) Empty() bool {
	return *e == (Event{})
}

func (e *Event) String() string {
	return fmt.Sprintf("%p - %#v", e, e)
}
