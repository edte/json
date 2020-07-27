// @program: json
// @author: edte
// @create: 2020-07-26 18:16
package value

import (
	"encoding/json"
	"errors"
)

type Json struct {
	raw interface{}
	err error
}

func NewJson(data []byte) (*Json, error) {
	j := new(Json)
	err := json.Unmarshal(data, &j.raw)
	return j, err
}

func (j *Json) error(msg string) {
	j.err = errors.New(msg)
}

func (j *Json) Error() error {
	return j.err
}

func (j *Json) Raw() interface{} {
	return j.raw
}

func (j *Json) Marshal() (string, error) {
	m, err := json.Marshal(j.raw)
	return string(m), err
}

func (j *Json) Get(path ...interface{}) (*Value, error) {
	v := NewValue(j.raw)
	for _, p := range path {
		switch p.(type) {
		case string:
			m, err := j.GetByKey(p.(string))
			if err != nil {
				return nil, err
			}
			v = m
			j.raw = v.data
		case int:
			m, err := j.GetByNumber(p.(int))
			if err != nil {
				return nil, err
			}
			v = m
			j.raw = v.data
		default:
			return nil, errors.New("path not valid")
		}
	}
	return v, nil
}

func (j *Json) GetByKey(path string) (*Value, error) {
	v := NewValue(j.raw)
	d, ok := (v.data).(map[string]interface{})
	if !ok {
		return &Value{}, errors.New("failed to assert map[string]interface{}")
	}
	if _, ok = d[path]; !ok {
		return &Value{}, errors.New("target key not exist")
	}
	v.data = d[path]
	return v, nil
}

func (j *Json) GetByKeys(path ...string) (*Value, error) {
	v := NewValue(j.raw)
	var err error

	for _, k := range path {
		v, err = j.GetByKey(k)
		if err != nil {
			return nil, err
		}
		j.raw = v.data
	}

	return v, nil
}

func (j *Json) GetByNumber(index int) (*Value, error) {
	v := NewValue(j.raw)

	d, ok := (v.data).([]interface{})
	if !ok {
		return &Value{}, errors.New("failed to assert []interface{}")
	}
	if len(d) <= index {
		return &Value{}, errors.New("out of index")
	}
	return NewValue(d[index]), nil
}

func (j *Json) GetByNumbers(index ...int) (Value, error) {
	for _, i := range index {
		v, err := j.GetByNumber(i)
		if err != nil {
			return Value{}, err
		}

		j.raw = v.data
	}
	return Value{j.raw}, nil
}

func (j *Json) IsExist(path ...string) bool {
	v := NewValue(j.raw)

	for _, k := range path {
		d, ok := (v.data).(map[string]interface{})
		if !ok {
			return false
		}
		if _, ok = d[k]; !ok {
			return false
		}
		v.data = d[k]
	}

	return true
}

func (j *Json) Set(value interface{}, path ...interface{}) error {
	for _, p := range path {
		switch p.(type) {
		case string:
			err := j.SetByKeys(value, []string{p.(string)})
			if err != nil {
				return err
			}
		case int:
			err := j.SetByNumber(p.(int), value)
			if err != nil {
				return err
			}
		default:
			return errors.New("path not valid")
		}
	}
	return nil
}

func (j *Json) SetByKeys(value interface{}, path []string) error {
	// todo: check value valid

	v := NewValue(j.raw)
	var tt []map[string]interface{}

	for _, k := range path {
		d, ok := (v.data).(map[string]interface{})
		if !ok {
			return errors.New("not object")
		}

		if _, ok = d[k]; !ok {
			return errors.New("key not exist")
		}

		v.data = d[k]
		tt = append(tt, d)
	}
	//  这步 nb
	tt[len(tt)-1][path[len(path)-1]] = value

	j.raw = tt[0]

	return nil
}

func (j *Json) SetByNumber(index int, value interface{}) error {
	v := NewValue(j.raw)
	d, ok := (v.data).([]interface{})
	if !ok {
		return errors.New("failed to assert []interface{}")
	}
	d[index] = value
	j.raw = d
	return nil
}

func (j *Json) Del(path ...string) error {
	// todo: add del
	return nil
}
