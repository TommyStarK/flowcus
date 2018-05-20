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

func newExploratory() *exploratory {
	return &exploratory{
		nil,
		&sync.WaitGroup{},
		0,
		NewFifo(),
		nil,
		make(chan *Input, 1),
		nil,
	}
}

type exploratory struct {
	Report
	*sync.WaitGroup
	once     uint64
	in       *Fifo
	_tFuncIn BoxIF
	cin      chan *Input
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
		log.Fatalln("Error: Run() can be called only once")
	} else if e.manager == nil {
		log.Fatalln("You must register at least one test. Test function must have the following signature: func(*Test, Input)")
	} else if e._tFuncIn == nil {
		log.Fatalln("You must register an input")
	}

	atomic.AddUint64(&e.once, 1)
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go e._tFuncIn(e.cin)

	e.Add(1)
	go func(sig chan os.Signal) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("recover:", r)
			}

			e.Done()
		}()

		LoopSingleChan(sig, e.cin, e.in)

		for e.in.Len() > 0 {
			e.manager.StartWorkers(e.in.Pop().(*Input))
		}
	}(sig)
	e.Wait()

	close(sig)
	e.Report = NewReport(e.manager)
}
