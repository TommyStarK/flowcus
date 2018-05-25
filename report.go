package flowcus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

var (
	reportProperty       = "--- %s "
	reportCase           = "====== case n°%d\n"
	reportCaseTestName   = "\t* caller: %s\n"
	reportCaseTestDetail = "\t	> %s: "
)

func exportAndAppendTest(t Test, results *[]TestExported, success, count *int) {
	*count++
	if !t.Failed() {
		*success++
	}
	*results = append(*results, TestExported{
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
			report.Duration += tests[i].duration
			exportAndAppendTest(*tests[i], &results, &success, &count)
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
	fmt.Printf("%s %s\n", fmt.Sprintf(reportProperty, "date"), b.Date)
	fmt.Printf("%s %s\n", fmt.Sprintf(reportProperty, "duration"), b.Duration.String())
	fmt.Printf("%s %g\n", fmt.Sprintf(reportProperty, "success"), b.Coverage)
	fmt.Printf("%s %d\n", fmt.Sprintf(reportProperty, "number"), b.Number)
	for i, c := range b.Cases {
		fmt.Printf(reportCase, i+1)
		for _, t := range c {
			fmt.Printf(reportCaseTestName, t.Caller)
			fmt.Printf("%s %s\n", fmt.Sprintf(reportCaseTestDetail, "duration"), t.Duration.String())
			fmt.Printf("%s %t\n", fmt.Sprintf(reportCaseTestDetail, "sucess"), t.Success)
			fmt.Printf("%s %t\n", fmt.Sprintf(reportCaseTestDetail, "finished"), t.Finished)
			fmt.Printf("%s %t\n", fmt.Sprintf(reportCaseTestDetail, "skipped"), t.Skipped)
		}
	}
}
