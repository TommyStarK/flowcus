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

func newExploratory() *exploratory {
	return &exploratory{
		nil,
		0,
		NewFifo(),
		nil,
		make(chan *Input, 1),
		&sync.WaitGroup{},
		nil,
	}
}

type exploratory struct {
	Report
	once     uint64
	in       *Fifo
	_tFuncIn BoxIF
	cin      chan *Input
	wg       *sync.WaitGroup
	manager  *exploratoryBoxTestsManager
}

func (e *exploratory) Input(fn BoxIF) {
	e._tFuncIn = fn
}

func (e *exploratory) RegisterTests(tests ...BoxETF) {
	if e.manager == nil {
		e.manager = NewExploratoryBoxTestsManager()
	}

	e.manager.SetTasks(tests...)
}

func (e *exploratory) ReportToCLI() {
	if e.Report == nil {
		log.Fatalln("Unexpected error occurred. Report is nil")
	}

	e.Report.ReportToCLI()
}

func (e *exploratory) ReportToJSON(filename string) error {
	if e.Report == nil {
		log.Fatalln("Unexpected error occurred. Report is nil")
	}

	return e.Report.ReportToJSON(filename)
}

func (e *exploratory) Run() {
	if once := atomic.LoadUint64(&e.once); once == 1 {
		log.Println(ISRUNNING)
		return
	} else if e.manager == nil {
		log.Fatalln(NO_TEST_SET, BOXETF)
	} else if e._tFuncIn == nil {
		log.Fatalln(NO_INPUT_SET)
	}

	atomic.AddUint64(&e.once, 1)
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go e._tFuncIn(e.cin)

	e.wg.Add(1)
	go func(sig chan os.Signal) {
		defer func() {
			if r := recover(); r != nil {
				switch r.(type) {
				case syscall.Signal:
					if r.(syscall.Signal) == syscall.SIGINT {
						log.Println(CTRLC)
					}
				default:
					panic(errors.New(fmt.Sprintf("[Flowcus] %s", r)))
				}
			}

			e.wg.Done()
		}()

		LoopSingleChan(sig, e.cin, e.in)

		for e.in.Len() > 0 {
			e.manager.StartWorkers(e.in.Pop().(*Input))
		}
	}(sigchan)

	e.wg.Wait()
	close(sigchan)
	e.Report = NewReport(e.manager)
}
