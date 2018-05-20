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
	BOXETF        string = "func(*Test, Input)"
	BOXLTF        string = "func(*Test, Input, Output)"
	BOXNLTF       string = "func(*Test, []Input, []Output)"
	CTRLC         string = "Flowcus interrupted by the user (ctrl+c)"
	ISRUNNING     string = "Box is already running or tests have already been run"
	NO_INPUT_SET  string = "You must register an input function. Input function must have the following signature: func(chan<- *Input)"
	NO_OUTPUT_SET string = "You must register an ouput function. Ouput function must have the following signature: func(chan<- *Ouput)"
	NO_TEST_SET   string = "You must register at least one test. Test function must have the following signature: "
)

type BoxIF func(chan<- *Input)
type BoxOF func(chan<- *Output)
type BoxETF func(*Test, Input)
type BoxLTF func(*Test, Input, Output)
type BoxNLTF func(*Test, []Input, []Output)

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
Loop:
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
				break Loop
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
