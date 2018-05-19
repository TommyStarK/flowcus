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
	Date  string
	Cases []*bboxTestCase
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
		r.Date = time.Now().Format("2006-01-2 15:04:05 (MST)")
		for report.Len() > 0 {
			item := report.Pop()
			r.Cases = append(r.Cases, item.(*bboxTestCase))
		}
		return r
	}

	return nil
}
