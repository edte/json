// @program: json
// @author: edte
// @create: 2020-07-24 14:57
package encode

import (
	"bytes"
	"reflect"
	"sync"
)

func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()
	err := e.marshal(v, encOpts{escapeHTML: true})
	if err == nil {
		return nil, err
	}

	buf := append([]byte(nil), e.Bytes()...)
	return buf, nil
}

type encodeState struct {
	bytes.Buffer
	scratch  [64]byte
	ptrLevel uint
	ptrSeen  map[interface{}]struct{}
}

func newEncodeState() *encodeState {
	return &encodeState{ptrSeen: make(map[interface{}]struct{})}
}

func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) {
	e.reflectValue(reflect.ValueOf(v), opts)
	return nil
}

func (e *encodeState) reflectValue(v reflect.Value, opts encOpts) {
	// todo:sdf
	valueEncoder(v)(e, v, opts)
}

type encoderFunc func(e *encodeState, v reflect.Value, opts encOpts)

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}

var encoderCache sync.Map

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	return nil
}

func invalidValueEncoder(e *encodeState, v reflect.Value, _ encOpts) {
	e.WriteString("null")
}

type encOpts struct {
	quoted     bool
	escapeHTML bool
}
