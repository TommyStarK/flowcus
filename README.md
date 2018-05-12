# Flowcus

Flowcus is a data flow testing tool.

## Context

As a devOps trainee in a company providing analytics services, my first assignment was to write a software to perform integration tests on the data collection process. The data goes through different stages of treatment, my goal was to test these steps as well as the entire process flow.
The real challenge was to correctly setup the test and being able to perform it, not really to write the test itself. So I have
decided to rewrite the 1st version of my software to make it generic and reusable.

## Contribution

Each Contribution is welcomed and encouraged. I do not claim to cover each use cases nor completely master the Golang. If you encounter a non sense or any trouble, you can open an issue and I will be happy to discuss about it :)

## Install

```bash
    go get github.com/TommyStarK/flowcus
```

## Test

```bash
    go test -v --cover ./...
```

## Usage

The only way i found to link the data produced and the data consumed is an ID, no matter its type.
Before going deeper in your test, you will have to be able to retrieve this ID from the data consumed.

N.B: I am trying to find another way to link data consumed and produced without relying on a ID.

```go
    package main

    import flow "github.com/TommyStarK/flowcus"

    func performTest(omap *flow.OrderedMap, data interface{}) (interface{}, error) {
        // Write your test here
        // Return whether the id or an error if your test failed

        // 1- Use data (the data consumed) to retrieve the ID.
        // 2- Use omap.Get(ID) to retrieve the flow of your test and access the data stored from your producer.
        //  The ordered map returns whether a *Flow{} or nil.
        // 3- Write your logic
        // 4- Returns whether the ID if everything went well or an error if one of the previous steps failed.
    }

    func main() {
        flowcus := flow.NewFlowcus()

        flowcus.Producer(func(com chan<- *flow.Event) {
            //  Your logic

            // An event has two attributes, Id and Data, both interfaces.
            // Id is mandatory, it is the key of Flowcus operation
            // Data is optional, it allows you to store some data that you might need during
            // your test.
            event := &flow.Event{
                Id: "this is an ID",
                Data: "this is some data",
            }

            // Register a new test
            com <- event

            // You must close the channel when you are done producing
            close(com)
        })

        flowcus.Consumer(func(com chan<- *flow.Revent) {
            // Your logic

            // A Revent has two attributes, both interfaces and both mandatory.
            // Data: the data consumed that you want to test.
            // Func: the function to execute to perform your test.
            // this function must have the following signature func(*flow.OrderedMap, interface{})(interface{}, error)
            revent := &flow.Revent{
                Data: "this is some data consumed on which i want to perform some tests",
                Func: performTest,
            }

            // Register the data consumed to perform your test.
            com <- revent

            // You must close the channel when you are done consuming
            close(com)
        })

        flowcus.Start()
        // Write results to a JSON file named "report.json"
        flowcus.ReportToJSON("report.json")
    }
```
