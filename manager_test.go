package flowcus

import (
	"reflect"
	"runtime"
	"testing"
)

func testExplore(t *Test, i Input) {
	t.Log("test expore")
}

func testExplorePanic(t *Test, i Input) {
	panic("explore panic")
}

func testLinear(t *Test, i Input, o Output) {
	t.Log("test linear")
}

func testLinearPanic(t *Test, i Input, o Output) {
	panic("linear panic")
}

func testNonLinear(t *Test, i []Input, o []Output) {
	t.Log("test nonlinear")
}

func testNonLinearPanic(t *Test, i []Input, o []Output) {
	panic("nonlinear panic")
}

func TestExploratoryBoxTestsManager(t *testing.T) {
	ebtm := NewExploratoryBoxTestsManager()

	if ebtm.cases == nil || ebtm.cases.Len() > 0 {
		t.Error("cases should not be nil, should of type *Fifo with a 0 length")
	}

	if ebtm.tests == nil || ebtm.tests.Len() > 0 {
		t.Error("tests should not be nil, should of type map[interface{}]interface{} with a 0 length")
	}

	if ebtm.wg == nil {
		t.Error("wg should not be nil and should be of type *sync.WaitGroup")
	}

	key := runtime.FuncForPC(reflect.ValueOf(testExplore).Pointer()).Name()
	ebtm.SetTasks(testExplore)
	ebtm.StartWorkers(&Input{})
	ebtm.tests.Delete(key)

	res := ebtm.cases.Pop()
	finished := res.(*exploratoryBoxTestCase).Results[0].finished
	success := !res.(*exploratoryBoxTestCase).Results[0].Failed()

	if !finished || !success {
		t.Error("test should be finished and successful")
	}

	ebtm.SetTasks(testExplorePanic)
	ebtm.StartWorkers(&Input{})

	res = ebtm.cases.Pop()
	finished = res.(*exploratoryBoxTestCase).Results[0].finished
	success = !res.(*exploratoryBoxTestCase).Results[0].Failed()

	if finished || success {
		t.Error("test should not be finished or successful")
	}
}

func TestLinearBoxTestsManager(t *testing.T) {
	lbtm := NewLinearBoxTestsManager()

	if lbtm.cases == nil || lbtm.cases.Len() > 0 {
		t.Error("cases should not be nil, should of type *Fifo with a 0 length")
	}

	if lbtm.tests == nil || lbtm.tests.Len() > 0 {
		t.Error("tests should not be nil, should of type map[interface{}]interface{} with a 0 length")
	}

	if lbtm.wg == nil {
		t.Error("wg should not be nil and should be of type *sync.WaitGroup")
	}

	key := runtime.FuncForPC(reflect.ValueOf(testLinear).Pointer()).Name()
	lbtm.SetTasks(testLinear)
	lbtm.StartWorkers(&Input{}, &Output{})
	lbtm.tests.Delete(key)

	res := lbtm.cases.Pop()
	finished := res.(*linearBoxTestCase).Results[0].finished
	success := !res.(*linearBoxTestCase).Results[0].Failed()

	if !finished || !success {
		t.Error("test should be finished and successful")
	}

	lbtm.SetTasks(testLinearPanic)
	lbtm.StartWorkers(&Input{}, &Output{})

	res = lbtm.cases.Pop()
	finished = res.(*linearBoxTestCase).Results[0].finished
	success = !res.(*linearBoxTestCase).Results[0].Failed()

	if finished || success {
		t.Error("test should not be finished or successful")
	}
}

func TestNonLinearBoxTestsManager(t *testing.T) {
	nlbtm := NewNonLinearBoxTestsManager()

	if nlbtm.cases == nil || nlbtm.cases.Len() > 0 {
		t.Error("cases should not be nil, should of type *Fifo with a 0 length")
	}

	if nlbtm.tests == nil || nlbtm.tests.Len() > 0 {
		t.Error("tests should not be nil, should of type map[interface{}]interface{} with a 0 length")
	}

	if nlbtm.wg == nil {
		t.Error("wg should not be nil and should be of type *sync.WaitGroup")
	}

	key := runtime.FuncForPC(reflect.ValueOf(testNonLinear).Pointer()).Name()
	nlbtm.SetTasks(testNonLinear)
	nlbtm.StartWorkers([]Input{Input{}, Input{}}, []Output{Output{}})
	nlbtm.tests.Delete(key)

	res := nlbtm.cases.Pop()
	finished := res.(*nonlinearBoxTestCase).Results[0].finished
	success := !res.(*nonlinearBoxTestCase).Results[0].Failed()

	if !finished || !success {
		t.Error("test should be finished and successful")
	}

	nlbtm.SetTasks(testNonLinearPanic)
	nlbtm.StartWorkers([]Input{Input{}, Input{}}, []Output{Output{}})

	res = nlbtm.cases.Pop()
	finished = res.(*nonlinearBoxTestCase).Results[0].finished
	success = !res.(*nonlinearBoxTestCase).Results[0].Failed()

	if finished || success {
		t.Error("test should not be finished or successful")
	}
}
