package flowcus

import (
	"os"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

const (
	_in int = iota
	_out
)

const (
	isRunning       = "Box is already running or tests have already been run"
	nilReport       = "Flowcus cannot synthesize a report"
	noInputFuncSet  = "You must register an input function. Input function must have the following signature: func(chan<- *Input)"
	noOutputFuncSet = "You must register an ouput function. Ouput function must have the following signature: func(chan<- *Ouput)"
	noTestSet       = "You must register at least one test. Test function must have the following signature: "
	sigINT          = "Flowcus interrupted by the user (ctrl+c)"
)

type (
	BoxIF   func(chan<- *Input)
	BoxOF   func(chan<- *Output)
	BoxETF  func(*Test, Input)
	BoxLTF  func(*Test, Input, Output)
	BoxNLTF func(*Test, []Input, []Output)
)

type Exploratory interface {
	Input(BoxIF)
	RegisterTests(...BoxETF)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

type Linear interface {
	Input(BoxIF)
	Output(BoxOF)
	RegisterTests(...BoxLTF)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

type NonLinear interface {
	Input(BoxIF)
	Output(BoxOF)
	RegisterTests(...BoxNLTF)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

func LoopSingleChan(sig chan os.Signal, cin chan *Input, in *Fifo) {
	for {
		select {
		case signal := <-sig:
			panic(signal)

		case input, open := <-cin:
			if open {
				if input != nil && !input.Empty() {
					in.Push(input)
				}
			} else {
				return
			}
		}
	}
}

func LoopDualChan(sig chan os.Signal, cin chan *Input, cout chan *Output, in, out *Fifo) {
	watcher := map[int]bool{_in: false, _out: false}

	for !watcher[_in] || !watcher[_out] {
		select {
		case signal := <-sig:
			panic(signal)

		case input, open := <-cin:
			if open {
				if input != nil && !input.Empty() {
					in.Push(input)
				}
			} else if !open && !watcher[_in] {
				watcher[_in] = true
			}

		case output, open := <-cout:
			if open {
				if output != nil && !output.Empty() {
					out.Push(output)
				}
			} else if !open && !watcher[_out] {
				watcher[_out] = true
			}
		}
	}
}
