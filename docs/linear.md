# Example

```go

    package main

    import . "github.com/TommyStarK/flowcus"

    var inputs []int = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    var outputs []int = []int{0, 1, 2, 3, 4, 0, 6, 7, 8, 9}

    func input(c chan<- *Input) {
        defer func() {
            close(c)
        }()

        for index := 0; index < len(inputs); index++ {
            c <- &Input{
                Data:     inputs[index],
                Expected: inputs[index],
            }
        }
    }

    func output(c chan<- *Output) {
        defer func() {
            close(c)
        }()

        for index := 0; index < len(outputs); index++ {
            c <- &Output{
                Data: outputs[index],
            }
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
        } else {
            t.Log("output is as expected")
        }
    }

    func main() {
        f := NewLinearBox()
        f.Input(input)
        f.Output(output)
        f.RegisterTests(test)
        f.Run()
        f.ReportToCLI()
        f.ReportToJSON("equivalency.json")
    }
```