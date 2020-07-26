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

func (j *Json) Get(key ...string) (*Value, error) {
	v := NewValue(j.raw)

	for _, k := range key {
		d, ok := (v.data).(map[string]interface{})
		if !ok {
			return &Value{}, errors.New("target key not exist")
		}
		if _, ok = d[k]; !ok {
			return &Value{}, errors.New("target key not exist")
		}
		v.data = d[k]
	}

	return v, nil
}

func (j *Json) IsExist(key ...string) bool {
	v := NewValue(j.raw)

	for _, k := range key {
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

func (j *Json) Set(value interface{}, key []string) error {
	v := NewValue(j.raw)
	var tt []map[string]interface{}

	for _, k := range key {
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
	tt[len(tt)-1][key[len(key)-1]] = value

	j.raw = tt[0]

	return nil
}

func (j *Json) Del() {

}
