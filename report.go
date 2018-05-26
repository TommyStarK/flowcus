package flowcus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	. "github.com/TommyStarK/flowcus/internal/decorator"
	. "github.com/TommyStarK/flowcus/internal/fifo"
)

type Report interface {
	ReportToCLI()
	ReportToJSON(string) error
}

func NewReport(manager interface{}) Report {
	var fifo *Fifo
	success, count := 0, 0
	report := new(BoxReport)
	cases := make([][]TestExported, 0)
	report.Date = time.Now().Format(FORMAT)

	switch manager.(type) {
	case *exploratoryBoxTestsManager:
		report.Box = "Exploratory Box"
		fifo = manager.(*exploratoryBoxTestsManager).cases

	case *linearBoxTestsManager:
		report.Box = "Linear Box"
		fifo = manager.(*linearBoxTestsManager).cases

	case *nonlinearBoxTestsManager:
		report.Box = "Non Linear Box"
		fifo = manager.(*nonlinearBoxTestsManager).cases

	default:
		return nil
	}

	for fifo.Len() > 0 {
		tests := fifo.Pop().([]*Test)
		results := make([]TestExported, 0)
		for i := 0; i < len(tests); i++ {
			count++
			if !tests[i].Failed() {
				success++
			}
			report.Duration += tests[i].duration
			results = append(results, TestExported{
				Caller:   tests[i].caller,
				Start:    tests[i].start,
				Duration: tests[i].duration,
				Finished: tests[i].finished,
				Skipped:  tests[i].Skipped(),
				Success:  !tests[i].Failed(),
				Errors:   tests[i].errors,
				Logs:     tests[i].logs,
			})
		}
		cases = append(cases, results)
	}

	report.Number = count
	report.Coverage = float64(success) / float64(count) * float64(100)
	report.Cases = cases
	return report
}

type BoxReport struct {
	Box      string
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    [][]TestExported
}

func (b *BoxReport) ReportToJSON(filename string) error {
	report, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, report, 0644)
}

func (b *BoxReport) ReportToCLI() {
	format := func(s string) string {
		res := []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
		for i := 0; i < len(s); i++ {
			res[i] = s[i]
		}

		return string(res)
	}

	fmt.Printf(
		"[%s] Tests took %s. %g%% of tests succeed for a total of %d tests performed (%s).\n",
		Colorize("purple", "Flowcus"),
		b.Duration.String(),
		b.Coverage,
		b.Number,
		b.Date,
	)
	for i, c := range b.Cases {
		fmt.Printf("--- case n°%s\n", Colorize("yellow", strconv.Itoa(i+1)))
		for _, t := range c {
			fmt.Printf("\t* caller: %s\n", t.Caller)
			fmt.Printf("%s %s\n", fmt.Sprintf("\t	> %s", format("duration:")), t.Duration.String())
			fmt.Printf("%s %s\n", fmt.Sprintf("\t	> %s", format("success:")), BoolToColorizedString(t.Success))
			fmt.Printf("%s %s\n", fmt.Sprintf("\t	> %s", format("finished:")), BoolToColorizedString(t.Finished))
			fmt.Printf("%s %s\n", fmt.Sprintf("\t	> %s", format("skipped:")), BoolToColorizedString(t.Skipped))
			for _, log := range t.Logs {
				f := "%s %s"
				if log[len(log)-1] != '\n' {
					f += "\n"
				}
				fmt.Printf(f, fmt.Sprintf("\t	> %s", format("log:")), log)
			}

			for _, err := range t.Errors {
				f := "%s %s"
				if err[len(err)-1] != '\n' {
					f += "\n"
				}
				fmt.Printf(f, fmt.Sprintf("\t	> %s", format("error:")), err)
			}
		}
	}
}
