package flowcus

import (
	"strings"
)

const (
	_in int = iota
	_out
)

type tFunc func(*Test, *Input, *Output)
type tFuncIn func(chan<- *Input)
type tFuncOut func(chan<- *Output)

type BlackBox interface {
	Input(tFuncIn)
	Output(tFuncOut)
	RegisterTests(...tFunc)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

func NewBlackBox(Type string) BlackBox {
	switch strings.ToLower(Type) {
	case "equivalency":
		return _equivalency()
	}

	return nil
}
