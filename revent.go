package flowcus

import (
	"fmt"
)

type Revent struct {
	Data interface{}
	Test interface{}
}

func (r *Revent) Empty() bool {
	return *r == (Revent{})
}

func (r *Revent) String() string {
	return fmt.Sprintf("%p - %#v", r, r)
}
