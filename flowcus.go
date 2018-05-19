package flowcus

const (
	VERSION float64 = 0.1
	FORMAT  string  = "2006-01-2 15:04:05 (MST)"
	_in     int     = iota
	_out
)

type tFuncIn func(chan<- *Input)
type tFuncOut func(chan<- *Output)
type tGBoxFunc func(*Test, Input)
type tBBoxFunc func(*Test, Input, Output)
