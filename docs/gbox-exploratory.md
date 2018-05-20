# Example

```go

    package main

    import . "github.com/TommyStarK/flowcus"

    func input(c chan<- *Input) {
        defer func() {
            close(c)
        }()

        c <- &Input{
            Data: map[int]int{2: 2},
        }

        c <- &Input{
            Data: true,
        }

        c <- &Input{
            Data: "data",
        }

        c <- &Input{
            Data: 42,
        }
    }

    func testA(t *Test, i Input) {
        switch i.Data.(type) {
        case int:
            t.Log("wanted type received")
        case string:
            t.Error("wanted int received string")
        case bool:
            t.Error("wanted int received bool")
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
        e := NewEploratoryBox()
        e.Input(input)
        e.RegisterTests(testA, testB, testC)
        e.Run()
        e.ReportToCLI()
        e.ReportToJSON("exploratory.json")
    }

```