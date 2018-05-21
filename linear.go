package flowcus

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

func newLinear() *linear {
	return &linear{
		nil,
		0,
		NewFifo(),
		NewFifo(),
		nil,
		nil,
		make(chan *Input, 1),
		make(chan *Output, 1),
		&sync.WaitGroup{},
		nil,
	}
}

type linear struct {
	Report
	once      uint64
	in        *Fifo
	out       *Fifo
	_tFuncIn  BoxIF
	_tFuncOut BoxOF
	cin       chan *Input
	cout      chan *Output
	wg        *sync.WaitGroup
	manager   *linearBoxTestsManager
}

func (l *linear) Input(fn BoxIF) {
	l._tFuncIn = fn
}

func (l *linear) Output(fn BoxOF) {
	l._tFuncOut = fn
}

func (l *linear) RegisterTests(tests ...BoxLTF) {
	if l.manager == nil {
		l.manager = NewLinearBoxTestsManager()
	}

	l.manager.SetTasks(tests...)
}

func (l *linear) ReportToCLI() {
	if l.Report == nil {
		log.Println(nilReport)
		return
	}

	l.Report.ReportToCLI()
}

func (l *linear) ReportToJSON(filename string) error {
	if l.Report == nil {
		return errors.New(nilReport)
	}

	return l.Report.ReportToJSON(filename)
}

func (l *linear) Run() {
	if once := atomic.LoadUint64(&l.once); once == 1 {
		log.Println(isRunning)
		return
	} else if l.manager == nil {
		log.Println(noTestSet, "func(*Test, Input, Output)")
		return
	} else if l._tFuncIn == nil {
		log.Println(noInputFuncSet)
		return
	} else if l._tFuncOut == nil {
		log.Println(noOutputFuncSet)
		return
	}

	atomic.AddUint64(&l.once, 1)
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go l._tFuncIn(l.cin)
	go l._tFuncOut(l.cout)

	l.wg.Add(1)
	go func(sig chan os.Signal) {
		defer func() {
			if r := recover(); r != nil {
				switch r.(type) {
				case syscall.Signal:
					if r.(syscall.Signal) == syscall.SIGINT {
						log.Println(sigTERM)
					}
				default:
					panic(errors.New(fmt.Sprintf("[Flowcus] %s", r)))
				}
			}

			l.wg.Done()
		}()

		LoopDualChan(sig, l.cin, l.cout, l.in, l.out)

		if l.out.Len() != l.in.Len() {
			log.Println("Different number of inputs/outputs. LinearBox exiting...")
			return
		}

		for l.in.Len() > 0 && l.out.Len() > 0 {
			l.manager.StartWorkers(l.in.Pop().(*Input), l.out.Pop().(*Output))
		}
	}(sigchan)

	l.wg.Wait()
	close(sigchan)
	l.Report = NewReport(l.manager)
}
