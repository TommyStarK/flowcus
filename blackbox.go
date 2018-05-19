package flowcus

import (
	"strings"
)

type BlackBox interface {
	Input(tFuncIn)
	Output(tFuncOut)
	RegisterTests(...tBBoxFunc)
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
