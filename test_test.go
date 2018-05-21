package flowcus

import (
	"sync"
	"testing"
)

func TestRawTest(t *testing.T) {
	raw := NewTest()

	if raw.caller != "" {
		t.Errorf("caller should be empty. caller: %s", raw.caller)
	}

	if raw.duration > 0 {
		t.Errorf("duration should be 0. duration: %d", raw.duration)
	}

	if raw.failed {
		t.Errorf("failed should be false. failed: %t", raw.failed)
	}

	if raw.finished {
		t.Errorf("finished should be false. finished: %t", raw.finished)
	}

	if raw.success {
		t.Errorf("success should be false. success: %t", raw.success)
	}

	if raw.skipped {
		t.Errorf("skipped should be false. skipped: %t", raw.skipped)
	}

	if len(raw.errors) > 0 {
		t.Errorf("len of errors should be 0. errors len: %d", len(raw.errors))
	}

	if len(raw.logs) > 0 {
		t.Errorf("len of logs should be 0. logs len: %d", len(raw.logs))
	}
}

func TestLog(t *testing.T) {
	raw := NewTest()

	raw.Log("test 1")
	raw.LogF("test %d", 2)

	if len(raw.logs) != 2 {
		t.Fatalf("two logs have been set, len of logs should be equal to 2. len %d", len(raw.logs))
	}

	if raw.logs[0] != "test 1\n" || raw.logs[1] != "test 2" {
		t.Error("logs corrupted")
	}
}

func TestFail(t *testing.T) {
	raw := NewTest()

	raw.Fail()

	if !raw.Failed() {
		t.Error("calling Fail() should set failed to true")
	}
}

func TestFailNow(t *testing.T) {
	raw := NewTest()
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
		}()
		raw.FailNow()
	}(&wg)
	wg.Wait()

	if !raw.Failed() || !raw.finished {
		t.Error("calling FailNow() should set both failed and finished to true")
	}
}

func TestFatalAndFatalF(t *testing.T) {
	raw1 := NewTest()
	raw2 := NewTest()

	var wg sync.WaitGroup

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
		}()
		raw1.Fatal("Fatal")
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
		}()
		raw2.FatalF("FatalF")
	}(&wg)

	wg.Wait()

	if !raw1.Failed() || !raw1.finished {
		t.Error("raw1 should have failed set to false and finished set to true")
	}

	if !raw2.Failed() || !raw2.finished {
		t.Error("raw2 should have failed set to false and finished set to true")
	}

	if raw1.errors[0] != "Fatal\n" || raw2.errors[0] != "FatalF" {
		t.Error("errors corrupted")
	}
}

func TestSkipSkipFAndSkipNow(t *testing.T) {
	raw := NewTest()

	raw.Skip("skip")

	if !raw.Skipped() {
		t.Error("skipped should be true after Skip()")
	}

	raw.skipped = false

	raw.SkipF("skip %s", "f")

	if !raw.Skipped() {
		t.Error("skipped should be true after SkipF()")
	}

	raw.skipped = false

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
		}()

		raw.SkipNow()
	}(&wg)
	wg.Wait()

	if !raw.Skipped() {
		t.Error("skipped should be true after SkipNow()")
	}
}

func TestErrorAndErrorF(t *testing.T) {
	raw := NewTest()

	raw.Error("toto")
	raw.ErrorF("titi")

	if len(raw.errors) != 2 {
		t.Fatalf("two errors have been set, len of errors should be equal to 2. len %d", len(raw.errors))
	}

	if raw.errors[0] != "toto\n" || raw.errors[1] != "titi" {
		t.Error("errors corrupted")
	}
}

func TestRawTestExported(t *testing.T) {
	raw := &TestExported{}

	if raw.Caller != "" {
		t.Error("caller should be empty")
	}

	if raw.Duration != 0 {
		t.Error("duration should be 0")
	}

	if len(raw.Errors) > 0 {
		t.Error("errors' len should be 0")
	}

	if len(raw.Logs) > 0 {
		t.Error("logs' len should be 0")
	}

	if raw.Finished {
		t.Error("finished should be false")
	}

	if raw.Success {
		t.Error("success should be false")
	}

	if raw.Skipped {
		t.Error("skipped should be false")
	}
}
