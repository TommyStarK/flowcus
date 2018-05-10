package main

import (
	"log"

	"github.com/TommyStarK/flowcus/flowcus"
)

func main() {
	fifo := flowcus.NewFifo()
	orderedMap := flowcus.NewOrderedMap()

	log.Printf("%+v", fifo)
	log.Printf("%+v", orderedMap)
}
