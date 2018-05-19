package flowcus

import (
	"os"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

const (
	_in int = iota
	_out
)

type tBoxIF func(chan<- *Input)
type tBoxOF func(chan<- *Output)
type tBoxSCTF func(*Test, Input)
type tBoxDCTF func(*Test, Input, Output)

type BoxSingleChan interface {
	Input(tBoxIF)
	RegisterTests(...tBoxSCTF)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

type BoxDualChan interface {
	Input(tBoxIF)
	Output(tBoxOF)
	RegisterTests(...tBoxDCTF)
	ReportToCLI()
	ReportToJSON(string) error
	Run()
}

func tLoopSingleChan(sig chan os.Signal, cin chan *Input, in *Fifo) {
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

func tLoopDualChan(sig chan os.Signal, cin chan *Input, cout chan *Output, in, out *Fifo) {
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
