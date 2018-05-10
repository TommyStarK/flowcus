package flowcus

import (
	"sync"
)

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		make([]interface{}, 0),
		make(map[interface{}]interface{}),
		&sync.Mutex{},
	}
}

type OrderedMap struct {
	keys  []interface{}
	mmap  map[interface{}]interface{}
	mutex *sync.Mutex
}

func (o *OrderedMap) Delete(key interface{}) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for i := len(o.keys) - 1; i >= 0; i-- {
		if o.keys[i] == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			delete(o.mmap, key)
			return true
		}
	}

	return false
}

func (o *OrderedMap) Get(key interface{}) interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if value, ok := o.mmap[key]; ok {
		return value
	}

	return nil
}

func (o *OrderedMap) Keys() []interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.keys
}

func (o *OrderedMap) Len() int {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return len(o.keys)
}

func (o *OrderedMap) Set(key interface{}, value interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if _, ok := o.mmap[key]; !ok {
		o.keys = append(o.keys, key)
		o.mmap[key] = value
		return
	}

	o.mmap[key] = value
}
