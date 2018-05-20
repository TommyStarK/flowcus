# Example

```go

    package main

    import . "github.com/TommyStarK/flowcus"

    func input(c chan<- *Input) {
        defer func() {
            close(c)
        }()

        c <- &Input{
            Data:  map[int]int{2: 2},
        }

        c <- &Input{
            Data:  true,
        }

        c <- &Input{
            Data:  "data",
        }

        c <- &Input{
            Data:  42,
        }
    }

    func testA(t *Test, i Input) {
        switch i.Data.(type) {
        case int:
            t.Log("good type")
        case string:
            t.Error("expected int received string")
        case bool:
            t.Error("exepected int received bool")
        default:
            t.Fail()
            t.SkipNow()
        }
    }

    func testB(t *Test, i Input) {
        t.Log("func B...")
    }

    func testC(t *Test, i Input) {
        t.Log("func C...")
    }

    func main() {
        e := NewExploratoryBox()
        e.Input(input)
        e.RegisterTests(testA, testB, testC)
        e.Run()
        e.ReportToJSON("exploratory.json")
    }

```