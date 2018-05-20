package flowcus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

func reportToJSON(filename string, data interface{}) error {
	report, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, report, 0644)
}

type Report interface {
	ReportToCLI()
	ReportToJSON(string) error
}

type boxSingleChanReportCase struct {
	Input   Input
	Results []TestExported
}

type boxSingleChanReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*boxSingleChanReportCase
}

func (b *boxSingleChanReport) ReportToCLI() {
	fmt.Printf("Flowcus: Report [%s]\n", b.Date)
	fmt.Printf("Tests took: %s ending with %g%% of success for a total of %d tests performed.\n", b.Duration.String(), b.Coverage, b.Number)
	for i, c := range b.Cases {
		fmt.Printf("==> Input n° %d\n", i)
		fmt.Printf("%+v\n==> Results\n", c.Input)
		for _, t := range c.Results {
			fmt.Printf("%+v", t)
		}
	}
}

func (b *boxSingleChanReport) ReportToJSON(filename string) error {
	return reportToJSON(filename, b)
}

type boxDualChanReportCase struct {
	Input   Input
	Output  Output
	Results []TestExported
}

type boxDualChanReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*boxDualChanReportCase
}

func (b *boxDualChanReport) ReportToCLI() {
	log.Println("Reporting to CLI...")
}

func (b *boxDualChanReport) ReportToJSON(filename string) error {
	return reportToJSON(filename, b)
}

// Find a proper way, duplicate code
func NewReport(Type string, report *Fifo) Report {
	switch Type {
	case "boxSingleChanReport":
		r := new(boxSingleChanReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for report.Len() > 0 {
			item := report.Pop()
			test := new(boxSingleChanReportCase)
			test.Input = item.(*boxSingleChanTestCase).Input

			for i := 0; i < len(item.(*boxSingleChanTestCase).Results); i++ {
				var t Test = *item.(*boxSingleChanTestCase).Results[i]

				count++
				if !t.Failed() {
					success++
				}
				r.Duration += t.duration
				test.Results = append(test.Results, TestExported{
					Caller:   t.caller,
					Start:    t.start,
					Duration: t.duration,
					Finished: t.finished,
					Skipped:  t.Skipped(),
					Success:  !t.Failed(),
					Errors:   t.errors,
					Logs:     t.logs,
				})
			}
			r.Cases = append(r.Cases, test)
		}

		r.Number = count
		r.Coverage = float64(success) / float64(count) * float64(100)
		return r

	case "boxDualChanReport":
		r := new(boxDualChanReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for report.Len() > 0 {
			item := report.Pop()
			test := new(boxDualChanReportCase)
			test.Input = item.(*boxDualChanTestCase).Input
			test.Output = item.(*boxDualChanTestCase).Output

			for i := 0; i < len(item.(*boxDualChanTestCase).Results); i++ {
				var t Test = *item.(*boxDualChanTestCase).Results[i]

				count++
				if !t.Failed() {
					success++
				}
				r.Duration += t.duration
				test.Results = append(test.Results, TestExported{
					Caller:   t.caller,
					Start:    t.start,
					Duration: t.duration,
					Finished: t.finished,
					Skipped:  t.Skipped(),
					Success:  !t.Failed(),
					Errors:   t.errors,
					Logs:     t.logs,
				})
			}
			r.Cases = append(r.Cases, test)
		}

		r.Number = count
		r.Coverage = float64(success) / float64(count) * float64(100)
		return r
	}

	return nil
}
