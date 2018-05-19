package reflect

import (
	"reflect"
	"runtime"
)

func GetEntityName(target interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(target).Pointer()).Name()
}
