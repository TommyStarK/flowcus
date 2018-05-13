package flowcus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fatih/color"
)

const (
	VERSION float64 = 0.1
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	purple = color.New(color.FgMagenta).SprintFunc()
)

func NewFlowcus() *Flowcus {
	return &Flowcus{
		NewFifo(),
		NewOrderedMap(),
		&sync.Mutex{},
		&sync.WaitGroup{},
		make(chan *Event, 1),
		make(chan *Revent, 1),
		nil,
		nil,
		nil,
		0,
	}
}

type Flowcus struct {
	jobs     *Fifo
	tests    *OrderedMap
	mutex    *sync.Mutex
	waitGrp  *sync.WaitGroup
	event    chan *Event
	revent   chan *Revent
	producer func(chan<- *Event)
	consumer func(chan<- *Revent)
	report   *Report
	once     uint64
}

func (f *Flowcus) synthesize() {
	report := &Report{
		Coverage: 0,
		Date:     time.Now().Format("2006-01-2 15:04:05 (MST)"),
		Duration: 0,
		Version:  VERSION,
	}

	success := 0
	for _, key := range f.tests.Keys() {
		if flow := f.tests.Get(key); flow != nil {
			test := &Test{
				Id:       key,
				Duration: flow.(*Flow).duration,
				Sample:   flow.(*Flow).sample,
				Success:  flow.(*Flow).success,
				Tester:   flow.(*Flow).tester,
			}

			if test.Success {
				success++
			}
			report.Duration += test.Duration
			report.Tests = append(report.Tests, test)
		}
	}

	report.Coverage = float64(success) / float64(f.tests.Len()) * float64(100)
	f.report = report
}

func (f *Flowcus) process() {
	if f.jobs.Len() > 0 {
		job := f.jobs.Pop()
		if job.(*Revent).Test == nil {
			log.Println("no func provided in Revent")
			return
		}

		switch job.(*Revent).Test.(type) {
		case func(*OrderedMap, interface{}) (interface{}, error):
			start := time.Now()
			id, err := job.(*Revent).Test.(func(*OrderedMap, interface{}) (interface{}, error))(f.tests, job.(*Revent).Data)
			if err != nil {
				log.Println("Error executing Test from Revent:", err)
				return
			}

			if test := f.tests.Get(id); test != nil {
				test.(*Flow).duration = time.Since(start)
				test.(*Flow).sample = job.(*Revent).Data
				test.(*Flow).success = true
				test.(*Flow).tester = runtime.FuncForPC(reflect.ValueOf(job.(*Revent).Test).Pointer()).Name()
			}

		default:
			log.Println("wrong type func in Revent")
		}

		return
	}
	<-time.After(50 * time.Millisecond)
}

func (f *Flowcus) flowcus(sig chan os.Signal) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}

		f.waitGrp.Done()
	}()

	run := true
	for run {
		select {
		case signal := <-sig:
			panic(signal)

		case event, open := <-f.event:
			if !open {
				if f.jobs.Len() > 0 {
					f.process()
				} else {
					<-time.After(50 * time.Millisecond)
				}
			}

			if event != nil && !event.Empty() {
				f.tests.Set(event.Id, &Flow{Data: event.Data, duration: 0, success: false, tester: ""})
			}

		case revent, open := <-f.revent:
			if !open {
				if f.jobs.Len() > 0 {
					f.process()
				} else {
					run = false
				}
			}

			if revent != nil && !revent.Empty() {
				f.jobs.Push(revent)
			}

		default:
			f.process()
		}
	}
}

func (f *Flowcus) Consumer(fn func(chan<- *Revent)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.consumer = fn
}

func (f *Flowcus) Producer(fn func(chan<- *Event)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.producer = fn
}

func (f *Flowcus) Report() *Report {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.report
}

func (f *Flowcus) ReportAsString() (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	report, err := json.Marshal(f.report)
	if err != nil {
		return "", err
	}

	return string(report), nil
}

func (f *Flowcus) ReportToCLI() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.report == nil {
		return
	}

	log.Printf("[%s] Tests took: %s. %g%% of %s, %g%% of %s for a total of %d tests performed.",
		purple("Flowcus"),
		f.report.Duration.String(),
		f.report.Coverage,
		green("success"),
		float64(100)-f.report.Coverage,
		red("failure"),
		f.tests.Len())
}

func (f *Flowcus) ReportToJSON(filename string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	report, err := json.Marshal(f.report)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, report, 0644)
}

func (f *Flowcus) Start() {
	if once := atomic.LoadUint64(&f.once); once == 1 {
		log.Fatalf("Error: Start() can be called only once")
	}

	if f.producer == nil {
		log.Fatalf("Error: Flowcus requires a producer. Exiting.")
	}

	if f.consumer == nil {
		log.Fatalf("Error: Flowcus requires a consumer. Exiting.")
	}

	atomic.AddUint64(&f.once, 1)
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go f.producer(f.event)
	go f.consumer(f.revent)

	f.waitGrp.Add(1)
	go f.flowcus(sig)
	f.waitGrp.Wait()

	close(sig)
	f.synthesize()
}
