package flowcus

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

func newEquivalency() *equivalency {
	return &equivalency{
		nil,
		nil,
		&sync.WaitGroup{},
		0,
		NewFifo(),
		NewFifo(),
		nil,
		nil,
		make(chan *Input, 1),
		make(chan *Output, 1),
		map[int]bool{_in: false, _out: false},
	}
}

type equivalency struct {
	Report
	*bboxManager
	*sync.WaitGroup
	once      uint64
	in        *Fifo
	out       *Fifo
	_tFuncIn  tFuncIn
	_tFuncOut tFuncOut
	cin       chan *Input
	cout      chan *Output
	watcher   map[int]bool
}

func (e *equivalency) loop(sig chan os.Signal) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}

		e.Done()
	}()

	for !e.watcher[_in] || !e.watcher[_out] {
		select {
		case signal := <-sig:
			panic(signal)

		case input, open := <-e.cin:
			if open {
				if input != nil && !input.Empty() {
					e.in.Push(input)
				}
			} else if !open && !e.watcher[_in] {
				e.watcher[_in] = true
			}

		case output, open := <-e.cout:
			if open {
				if output != nil && !output.Empty() {
					e.out.Push(output)
				}
			} else if !open && !e.watcher[_out] {
				e.watcher[_out] = true
			}
		}
	}

	if e.out.Len() != e.in.Len() {
		log.Println("Number of inputs different from the number of outputs. An equivalence can not be made between two lists containing a different number of elements")
	}

	for e.in.Len() > 0 || e.out.Len() > 0 {
		e.StartWorkers(e.in.Pop().(*Input), e.out.Pop().(*Output))
	}
}

func (e *equivalency) Input(fn tFuncIn) {
	e._tFuncIn = fn
}

func (e *equivalency) Output(fn tFuncOut) {
	e._tFuncOut = fn
}

func (e *equivalency) RegisterTests(tests ...tBBoxFunc) {
	if e.bboxManager == nil {
		e.bboxManager = newBBoxManager()
	}

	e.SetTasks(tests...)
}

func (e *equivalency) ReportToCLI() {
	if e.Report == nil {
		log.Fatalln("Unexpected error occurred. Report is nil")
	}

	e.Report.ReportToCLI()
}

func (e *equivalency) ReportToJSON(filename string) error {
	if e.Report == nil {
		log.Fatalln("Unexpected error occurred. Report is nil")
	}

	return e.Report.ReportToJSON(filename)
}

func (e *equivalency) Run() {
	if once := atomic.LoadUint64(&e.once); once == 1 {
		log.Fatalln("Error: Run() can be called only once")
	} else if e.bboxManager == nil {
		log.Fatalln("You must register at least one test. Test function must have the following signature: func(*Test, *Input, *Output)")
	} else if e._tFuncIn == nil {
		log.Fatalln("You must register an input")
	} else if e._tFuncOut == nil {
		log.Fatalln("You must register an output")
	}

	atomic.AddUint64(&e.once, 1)
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go e._tFuncIn(e.cin)
	go e._tFuncOut(e.cout)

	e.Add(1)
	go e.loop(sig)
	e.Wait()

	close(sig)
	e.Report = newReport("bboxReport", e.bboxManager.Fifo)
}
