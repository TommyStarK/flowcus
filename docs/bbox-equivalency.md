# Example

```go

    package main

    import (
        "log"
        "time"

        . "github.com/TommyStarK/flowcus"
    )

    func input(c chan<- *Input) {
        defer func() {
            close(c)
        }()

        log.Println("sending input")

        for index := 0; index < 10; index++ {
            c <- &Input{
                Data:     index,
                Expected: index,
                Label:    "should succeed",
            }
        }

        c <- &Input{
            Data:     12,
            Expected: 12,
            Label:    "should fail",
        }
    }

    func output(c chan<- *Output) {
        defer func() {
            close(c)
        }()

        log.Println("sending output")
        time.Sleep(2 * time.Second)
        for index := 0; index < 10; index++ {
            time.Sleep(300 * time.Millisecond)
            c <- &Output{
                Data: index,
            }
        }

        c <- &Output{
            Data: 11,
        }
    }

    func test(t *Test, i Input, o Output) {
        if _, ok := i.Data.(int); !ok {
            t.Fatal("input should be of type int")
        }

        if _, ok := i.Expected.(int); !ok {
            t.Fatal("expected should be of type int")
        }

        if _, ok := o.Data.(int); !ok {
            t.Fatal("output should be of type int")
        }

        if i.Expected.(int) != o.Data.(int) {
            t.Error("expected and ouput should be equal")
        }
    }

    func main() {
        f := NewBoxDualChan("equivalency")
        f.Input(input)
        f.Output(output)
        f.RegisterTests(test)
        f.Run()
        f.ReportToJSON("equivalency.json")
    }
```