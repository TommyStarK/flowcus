# Example

```go
    package main

    import (
        "log"

        . "github.com/TommyStarK/flowcus"
    )

    func input(c chan<- *Input) {
        defer func() {
            close(c)
        }()

        log.Println("sending input")

        for index := 0; index < 10; index++ {
            c <- &Input{
                Data: index,
            }
        }
    }

    func output(c chan<- *Output) {
        defer func() {
            close(c)
        }()

        c <- &Output{
            Data: 10,
        }
    }

    func test(t *Test, inputs []Input, outputs []Output) {
        if outputs[0].Data.(int) != len(inputs) {
            t.Fatal("output should be equal to the total of inputs received")
        }

        t.Log("Success")
    }

    func main() {
        f := NewNonLinearBox()
        f.Input(input)
        f.Output(output)
        f.RegisterTests(test)
        f.Run()
        f.ReportToJSON("nonlinear.json")
    }
```