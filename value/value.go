// @program: json
// @author: edte
// @create: 2020-07-26 14:47
package value

import (
	"errors"
	"reflect"
)

type value interface {
	ToString() (string, error)
	ToBool() (bool, error)
	ToNull() (null, error)
	ToArray() (bool, error)
	ToObject() (bool, error)
	ToInt() (int, error)
	ToInt8() (int8, error)
	ToInt16() (int16, error)
	ToInt32() (int32, error)
	ToInt64() (int64, error)
	ToUint() (uint, error)
	ToUint8() (uint8, error)
	ToUint16() (uint16, error)
	ToUint32() (uint32, error)
	ToUint64() (uint64, error)
	ToFloat32() (float32, error)
	ToFloat64() (float64, error)
}

type null int

type Value struct {
	data interface{}
}

func NewValue(data interface{}) *Value {
	return &Value{data: data}
}

func (v *Value) Data() interface{} {
	return v.data
}

func (v *Value) Type() reflect.Type {
	return reflect.TypeOf(v.data)
}

func (v *Value) ToNull() (null, error) {
	if d, ok := (v.data).(null); ok {
		return d, nil
	}

	return -1, errors.New("not float32 float32")
}

func (v *Value) ToArray() ([]interface{}, error) {
	if d, ok := (v.data).([]interface{}); ok {
		return d, nil
	}

	return nil, errors.New("not []interface{} array")
}

func (v *Value) ToObject() (map[string]interface{}, error) {
	if d, ok := (v.data).(map[string]interface{}); ok {
		return d, nil
	}

	return nil, errors.New("not map[string]interface{} float32")
}

func (v *Value) ToFloat32() (float32, error) {
	if d, ok := (v.data).(float32); ok {
		return d, nil
	}

	return 0, errors.New("not float32 float32")
}

func (v *Value) ToFloat64() (float64, error) {
	if d, ok := (v.data).(float64); ok {
		return d, nil
	}

	return 0, errors.New("not float64 float64")
}

func (v *Value) ToInt() (int, error) {
	if d, ok := (v.data).(int); ok {
		return d, nil
	}

	return 0, errors.New("not int int")

}

func (v *Value) ToInt8() (int8, error) {
	if d, ok := (v.data).(int8); ok {
		return d, nil
	}

	return 0, errors.New("not int8 type")
}

func (v *Value) ToInt16() (int16, error) {
	if d, ok := (v.data).(int16); ok {
		return d, nil
	}

	return 0, errors.New("not int16 type")
}

func (v *Value) ToInt32() (int32, error) {
	if d, ok := (v.data).(int32); ok {
		return d, nil
	}

	return 0, errors.New("not int32 type")
}

func (v *Value) ToInt64() (int64, error) {
	if d, ok := (v.data).(int64); ok {
		return d, nil
	}

	return 0, errors.New("not int64 type")
}

func (v *Value) ToUint() (uint, error) {
	if d, ok := (v.data).(uint); ok {
		return d, nil
	}

	return 0, errors.New("not uint type")
}

func (v *Value) ToUint8() (uint8, error) {
	if d, ok := (v.data).(uint8); ok {
		return d, nil
	}

	return 0, errors.New("not uint8 type")
}

func (v *Value) ToUint16() (uint16, error) {
	if d, ok := (v.data).(uint16); ok {
		return d, nil
	}

	return 0, errors.New("not uint16 type")
}

func (v *Value) ToUint32() (uint32, error) {
	if d, ok := (v.data).(uint32); ok {
		return d, nil
	}

	return 0, errors.New("not uint32 type")
}

func (v *Value) ToUint64() (uint64, error) {
	if d, ok := (v.data).(uint64); ok {
		return d, nil
	}

	return 0, errors.New("not uint64 type")
}

func (v *Value) ToString() (string, error) {
	if d, ok := (v.data).(string); ok {
		return d, nil
	}

	return "", errors.New("not string type")
}

func (v *Value) ToBool() (bool, error) {
	if d, ok := (v.data).(bool); ok {
		return d, nil
	}

	return false, errors.New("not bool type")
}
