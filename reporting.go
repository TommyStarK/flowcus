package flowcus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	. "github.com/TommyStarK/flowcus/internal/fifo"
)

type Report interface {
	ReportToCLI()
	ReportToJSON(string) error
}

type bboxReport struct {
	Date     string
	Duration time.Duration
	Coverage float64
	Number   int
	Cases    []*bboxTestCase
}

func (b *bboxReport) ReportToCLI() {
	log.Println("Reporting to CLI...")
}

func (b *bboxReport) ReportToJSON(filename string) error {
	log.Println("Reporting to JSON...")

	report, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, report, 0644)
}

func newReport(Type string, report *Fifo) Report {
	switch Type {
	case "bboxReport":
		r := new(bboxReport)
		success, count := 0, 0
		r.Date = time.Now().Format(FORMAT)

		for report.Len() > 0 {
			item := report.Pop()
			for i := 0; i < len(item.(*bboxTestCase).Results); i++ {
				item.(*bboxTestCase).Results[i].Success = !item.(*bboxTestCase).Results[i].Failed()
				if item.(*bboxTestCase).Results[i].Success {
					success++
				}
				r.Duration += item.(*bboxTestCase).Results[i].Duration
				count++
			}
			r.Cases = append(r.Cases, item.(*bboxTestCase))
		}

		r.Number = count
		r.Coverage = float64(success) / float64(count) * float64(100)
		return r
	}

	return nil
}
