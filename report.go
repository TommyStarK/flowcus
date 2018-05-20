package flowcus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func reportToJSON(filename string, data interface{}) error {
	report, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, report, 0644)
}

// Find a proper way, duplicate code
func NewReport(manager interface{}) Report {
	switch manager.(type) {
	case *exploratoryBoxTestsManager:
		r := new(exploratoryBoxReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for manager.(*exploratoryBoxTestsManager).cases.Len() > 0 {
			item := manager.(*exploratoryBoxTestsManager).cases.Pop()
			test := new(exploratoryBoxReportCase)
			test.Input = item.(*exploratoryBoxTestCase).Input

			for i := 0; i < len(item.(*exploratoryBoxTestCase).Results); i++ {
				var t Test = *item.(*exploratoryBoxTestCase).Results[i]

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

	case *linearBoxTestsManager:
		r := new(linearBoxReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for manager.(*linearBoxTestsManager).cases.Len() > 0 {
			item := manager.(*linearBoxTestsManager).cases.Pop()
			test := new(linearBoxReportCase)
			test.Input = item.(*linearBoxTestCase).Input
			test.Output = item.(*linearBoxTestCase).Output

			for i := 0; i < len(item.(*linearBoxTestCase).Results); i++ {
				var t Test = *item.(*linearBoxTestCase).Results[i]

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

	case *nonlinearBoxTestsManager:
		r := new(nonLinearBoxReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for manager.(*nonlinearBoxTestsManager).cases.Len() > 0 {
			item := manager.(*nonlinearBoxTestsManager).cases.Pop()
			test := new(nonLinearBoxReportCase)
			test.Inputs = item.(*nonlinearBoxTestCase).Inputs
			test.Outputs = item.(*nonlinearBoxTestCase).Outputs

			for i := 0; i < len(item.(*nonlinearBoxTestCase).Results); i++ {
				var t Test = *item.(*nonlinearBoxTestCase).Results[i]

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

type Report interface {
	ReportToCLI()
	ReportToJSON(string) error
}

//
// Exploratory Box Report
//
type exploratoryBoxReportCase struct {
	Input   Input
	Results []TestExported
}

type exploratoryBoxReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*exploratoryBoxReportCase
}

func (b *exploratoryBoxReport) ReportToCLI() {
	fmt.Printf("Flowcus: Report [%s]\n", b.Date)
	fmt.Printf("Tests took: %s ending with %g%% of success for a total of %d tests performed.\n", b.Duration.String(), b.Coverage, b.Number)
	for i, c := range b.Cases {
		fmt.Printf("==> Input n° %d\n", i)
		fmt.Printf("%+v\n==> Results\n", c.Input)
		for _, t := range c.Results {
			fmt.Printf("%+v", t)
		}
		fmt.Printf("\n")
	}
}

func (b *exploratoryBoxReport) ReportToJSON(filename string) error {
	return reportToJSON(filename, b)
}

//
// Linear Box Report
//
type linearBoxReportCase struct {
	Input   Input
	Output  Output
	Results []TestExported
}

type linearBoxReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*linearBoxReportCase
}

func (b *linearBoxReport) ReportToCLI() {
	log.Println("Reporting to CLI...")
}

func (b *linearBoxReport) ReportToJSON(filename string) error {
	return reportToJSON(filename, b)
}

//
// Non Linear Box Report
//
type nonLinearBoxReportCase struct {
	Inputs  []Input
	Outputs []Output
	Results []TestExported
}

type nonLinearBoxReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*nonLinearBoxReportCase
}

func (n *nonLinearBoxReport) ReportToCLI() {
	log.Println("Reporting to CLI...")
}

func (n *nonLinearBoxReport) ReportToJSON(filename string) error {
	return reportToJSON(filename, n)
}
