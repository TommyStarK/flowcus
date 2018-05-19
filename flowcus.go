package flowcus

const (
	VERSION float64 = 0.1
	_in     int     = iota
	_out
)

type tFuncIn func(chan<- *Input)
type tFuncOut func(chan<- *Output)
type tGBoxFunc func(*Test, Input)
type tBBoxFunc func(*Test, Input, Output)
