package main

import (
	"log"
	"time"

	flow "github.com/TommyStarK/flowcus/flowcus"
)

func main() {
	fifo := flow.NewFifo()
	orderedMap := flow.NewOrderedMap()
	flowcus := flow.NewFlowcus()

	log.Printf("%+v", fifo)
	log.Printf("%+v", orderedMap)
	log.Printf("%+v\n\n\n", flowcus)

	flowcus.Producer(func(com chan<- *flow.Event) {
		log.Println("producer waiting for 1 sec")
		time.Sleep(1 * time.Second)
		log.Println("producing 2 events")
		event1 := &flow.Event{Id: 1}
		event2 := &flow.Event{Id: "this is an id"}
		com <- event1
		com <- event2
		close(com)
		// wg.Done()
		log.Println("Producer exiting")
	})

	flowcus.Consumer(func(com chan<- *flow.Revent) {
		log.Println("consumer waiting for 5 sec")
		time.Sleep(5 * time.Second)
		log.Println("consuming two revents")
		revent1 := &flow.Revent{Data: "1 data from consumer"}
		revent2 := &flow.Revent{Data: "2 data from consumer"}
		com <- revent1
		com <- revent2
		log.Println("consumer waiting for 5 sec")
		time.Sleep(5 * time.Second)
		close(com)
		// wg.Done()
		log.Println("Consumer exiting")
	})

	flowcus.Start()
}
