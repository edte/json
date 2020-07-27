// @program: json
// @author: edte
// @create: 2020-07-27 16:06
package encode

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

// 这个文件实现了 json 的序列化

// Marshal 传入一个 go 的数据类型，返回对应的 byte 流
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)
	return buf, nil
}

type encoderFunc func(e *encodeState, v reflect.Value)

type jsonError struct {
	error
}

func (e *encodeState) error(err error) {
	panic(jsonError{err})
}

// encodeState 是总的 encoder
type encodeState struct {
	// 这是 bytes buffer，encode 得到的 byte 就存这里面
	bytes.Buffer
}

// newEncodeState return a new encodeState
func newEncodeState() *encodeState {
	return &encodeState{}
}

// marshal 开始反射 value
func (e *encodeState) marshal(v interface{}) (err error) {
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

// reflectValue 比较重要，利用了 closure 的性质，把 encoder 作为返回值
// 这里起了一个中间站的作用，先反射拿到 encoder，再使用这个 encoder
func (e *encodeState) reflectValue(v reflect.Value) {
	// get encoder
	encoder := valueEncoder(v)
	// call encoder
	encoder(e, v)
}

// valueEncoder 判断 value 是否有效，然后开始调用 type encoder
func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}

// typeEncoder call type encoder
func typeEncoder(t reflect.Type) encoderFunc {
	return newTypeEncoder(t)
}

// newTypeEncoder 返回对应 type 的 encoder
func newTypeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.String:
		return stringEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.Bool:
		return boolEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Ptr:
		return ptrEncoder
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Slice, reflect.Array:
		return arrayOrSliceEncoder
	default:
		return unsupportedTypeEncoder
	}
}

// invalidValueEncoder 对应于无法解析的 value
func invalidValueEncoder(e *encodeState, v reflect.Value) {
	e.WriteString("invalidValueEncoder")
}

// ptrEncoder 用于指针的 encode
func ptrEncoder(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// 使用 reflect.Elem() 方法可以获得 ptr 指向的 type，然后调用所有的 encoder，当做 interface 即可
	e.reflectValue(v.Elem())
}

// arrayOrSliceEncoder 用于 slice 和 array 的 encode
func arrayOrSliceEncoder(e *encodeState, v reflect.Value) {
	e.WriteByte('[')
	// get slice length
	n := v.Len()
	// 遍历，对每个 value 当作 interface 来 encode
	for i := 0; i < n; i++ {
		e.reflectValue(v.Index(i))
		if i < n-1 {
			e.WriteByte(',')
		}
	}
	e.WriteByte(']')
}

// newMapEncoder 用于 map 格式的判断，然后选对应的 encode
func newMapEncoder(t reflect.Type) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String:
		return mapEncoder
	default:
		return invalidValueEncoder
	}
}

// mapEncoder encode map
func mapEncoder(e *encodeState, v reflect.Value) {
	e.WriteByte('{')

	// 获取 map key slice
	keys := v.MapKeys()

	// 遍历 key slice，key 直接 写入，主要是 value ，然后 value 当作 interface encode
	for i, key := range keys {
		e.WriteByte('"')
		e.WriteString(key.String())
		e.WriteByte('"')
		e.WriteByte(':')
		e.reflectValue(v.MapIndex(v.MapKeys()[i]))
		if i < len(keys)-1 {
			e.WriteByte(',')
		}
	}

	e.WriteByte('}')
}

// interfaceEncoder encode interface
func interfaceEncoder(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// 获得 interface 的动态类型
	e.reflectValue(v.Elem())
}

// stringEncoder encode string
func stringEncoder(e *encodeState, v reflect.Value) {
	e.WriteByte('"')
	e.WriteString(v.String())
	e.WriteByte('"')
}

// intEncoder encode int
func intEncoder(e *encodeState, v reflect.Value) {
	appendInt := strconv.AppendInt([]byte(nil), v.Int(), 10)
	e.Write(appendInt)
}

// uintEncoder encode uint
func uintEncoder(e *encodeState, v reflect.Value) {
	appendInt := strconv.AppendUint([]byte(nil), v.Uint(), 10)
	e.Write(appendInt)
}

// boolEncoder encode bool
func boolEncoder(e *encodeState, v reflect.Value) {
	if v.Bool() {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
}

// float32Encoder encode float32
func float32Encoder(e *encodeState, v reflect.Value) {
	b := strconv.AppendFloat([]byte(nil), v.Float(), 'f', -1, 32)
	e.Write(b)
}

// float64Encoder encode float64
func float64Encoder(e *encodeState, v reflect.Value) {
	b := strconv.AppendFloat([]byte(nil), v.Float(), 'f', -1, 64)
	e.Write(b)
}

// unsupportedTypeEncoder encode unsupported type
func unsupportedTypeEncoder(e *encodeState, v reflect.Value) {
	e.error(errors.New("unsupportedTypeEncoder"))
}

// newStructEncoder encode struct
func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{
		// 遍历存入 field
		field: typeFields(t),
	}
	return se.encode
}

// filed 表示 一个 filed
type field struct {
	name      string      // field name
	nameBytes []byte      // field name bytes
	encoder   encoderFunc // field type encoder
	value     reflect.Value
}

// structEncoder 表示一个 struct
type structEncoder struct {
	field []field
}

// typeFields 遍历存入 field
func typeFields(t reflect.Type) []field {
	fields := []field{}

	for j := 0; j < t.NumField(); j++ {
		fields = append(fields, field{
			name:      t.Field(j).Name,
			nameBytes: []byte(t.Field(j).Name),
			encoder:   typeEncoder(t.Field(j).Type),
			value:     reflect.ValueOf(t.Field(j)),
		})
	}

	return fields
}

// encode encode struct
func (se structEncoder) encode(e *encodeState, v reflect.Value) {
	next := byte('{')

	// 遍历 field，然后调用对应的 encoder 进行 encode
	for i := range se.field {
		f := &se.field[i]

		fv := v
		fv = fv.Field(i)

		e.WriteByte(next)
		e.WriteByte('"')
		e.Write(f.nameBytes)
		e.WriteByte('"')
		e.WriteByte(':')
		next = ','
		// call encoder
		f.encoder(e, fv)
	}

	if next == '{' {
		e.WriteString("{}")
	} else {
		e.WriteByte('}')
	}
}
