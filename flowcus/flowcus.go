package flowcus

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	}
}

type Flowcus struct {
	jobs     *Fifo
	tests    *OrderedMap
	mutex    *sync.Mutex
	wait     *sync.WaitGroup
	event    chan *Event
	revent   chan *Revent
	producer func(chan<- *Event)
	consumer func(chan<- *Revent)
}

func (f *Flowcus) synthetize() {
	log.Println(f.tests, "Flowcus exiting")
}

func (f *Flowcus) process() {
	if f.jobs.Len() > 0 {
		job := f.jobs.Pop()
		log.Printf("job: [%+v]\n", job)
		// TODO: inject user code here to perform their test
		// start := time.Now()
		// id, err := func (*OrderedMap, interface) (interface, error)
		// time.Since(start)
		// retrieve test struct by id
		// store result, err, etc
	} else {
		<-time.After(50 * time.Millisecond)
	}
}

func (f *Flowcus) flowcus(sig chan os.Signal) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}

		f.wait.Done()
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
				// TODO: IMPLEMENT STRUCT TEST
				f.tests.Set(event.Id, "TODO: IMPLEMENT STRUCT TEST")
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

func (f *Flowcus) Start() {
	if f.producer == nil || f.consumer == nil {
		log.Fatalf("Error: Flowcus requires a consumer and a producer.")
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go f.producer(f.event)
	go f.consumer(f.revent)

	f.wait.Add(1)
	go f.flowcus(sig)
	f.wait.Wait()

	close(sig)
	f.synthetize()
}
