package flowcus

import "testing"

func TestReport(t *testing.T) {
	test1 := &Test{}
	test2 := &Test{}
	report := &Report{}

	if !test1.Empty() || !test2.Empty() {
		t.Errorf("test1 and test2 should be empty")
	}

	if !report.Empty() {
		t.Errorf("report should be empty")
	}

	test1.Id = 1
	test1.Duration = 3
	test1.Sample = "sample"
	test1.Success = false
	test1.Tester = "main.test&"
	test2.Id = "2"
	test2.Duration = 23426
	test2.Sample = []byte("toto")
	test2.Success = true
	test2.Tester = "main.test2"

	if test1.Empty() || test2.Empty() {
		t.Errorf("test1 and test2 should not be empty")
	}

	report.Tests = append(report.Tests, test1)
	report.Tests = append(report.Tests, test2)

	if report.Empty() {
		t.Errorf("report should not be empty")
	}

}
