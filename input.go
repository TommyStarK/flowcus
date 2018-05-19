package flowcus

import "fmt"

type Input struct {
	Data     interface{}
	Expected interface{}
	Id       interface{}
	Label    string
}

func (i *Input) Empty() bool {
	return *i == (Input{})
}

func (i *Input) String() string {
	return fmt.Sprintf("%p - %#v", i, i)
}
