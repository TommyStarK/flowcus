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

func newNonLinear() *nonlinear {
	return &nonlinear{
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

type nonlinear struct {
	Report
	once      uint64
	in        *Fifo
	out       *Fifo
	_tFuncIn  BoxIF
	_tFuncOut BoxOF
	cin       chan *Input
	cout      chan *Output
	wg        *sync.WaitGroup
	manager   *nonlinearBoxTestsManager
}

func (n *nonlinear) Input(fn BoxIF) {
	n._tFuncIn = fn
}

func (n *nonlinear) Output(fn BoxOF) {
	n._tFuncOut = fn
}

func (n *nonlinear) RegisterTests(tests ...BoxNLTF) {
	if n.manager == nil {
		n.manager = NewNonLinearBoxTestsManager()
	}

	n.manager.SetTasks(tests...)
}

func (n *nonlinear) ReportToCLI() {
	if n.Report == nil {
		log.Println(nilReport)
		return
	}

	n.Report.ReportToCLI()
}

func (n *nonlinear) ReportToJSON(filename string) error {
	if n.Report == nil {
		return errors.New(nilReport)
	}

	return n.Report.ReportToJSON(filename)
}

func (n *nonlinear) Run() {
	if once := atomic.LoadUint64(&n.once); once == 1 {
		log.Println(isRunning)
		return
	} else if n.manager == nil {
		log.Println(noTestSet, "func(*Test, []Input, []Output)")
		return
	} else if n._tFuncIn == nil {
		log.Println(noInputFuncSet)
		return
	} else if n._tFuncOut == nil {
		log.Println(noOutputFuncSet)
		return
	}

	atomic.AddUint64(&n.once, 1)
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go n._tFuncIn(n.cin)
	go n._tFuncOut(n.cout)

	n.wg.Add(1)
	go func(sig chan os.Signal) {
		defer func() {
			if r := recover(); r != nil {
				switch r.(type) {
				case syscall.Signal:
					if r.(syscall.Signal) == syscall.SIGINT {
						log.Printf(sigINT)
					}
				default:
					panic(errors.New(fmt.Sprintf("[Flowcus] %s", r)))
				}
			}

			n.wg.Done()
		}()

		var inputs []Input
		var outputs []Output

		LoopDualChan(sig, n.cin, n.cout, n.in, n.out)

		for n.in.Len() > 0 {
			inputs = append(inputs, *n.in.Pop().(*Input))
		}

		for n.out.Len() > 0 {
			outputs = append(outputs, *n.out.Pop().(*Output))
		}

		n.manager.StartWorkers(inputs, outputs)
	}(sigchan)

	n.wg.Wait()
	close(sigchan)
	n.Report = NewReport(n.manager)
}
