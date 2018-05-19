package flowcus

import "fmt"

type Output struct {
	Data interface{}
}

func (o *Output) Empty() bool {
	return *o == (Output{})
}

func (o *Output) String() string {
	return fmt.Sprintf("%p - %#v", o, o)
}
