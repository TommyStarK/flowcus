package flowcus

import "strings"

const (
	FORMAT  string  = "2006-01-2 15:04:05 (MST)"
	VERSION float64 = 0.1
)

func NewBoxDualChan(Type string) BoxDualChan {
	switch strings.ToLower(Type) {
	case "equivalency":
		return newEquivalency()
	}

	return nil
}

func NewBoxSingleChan(Type string) BoxSingleChan {
	switch strings.ToLower(Type) {
	case "exploratory":
		return newExploratory()
	}

	return nil
}
