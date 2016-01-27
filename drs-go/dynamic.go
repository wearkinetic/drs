package drs

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ironbay/delta/encoding"
)

type Dynamic map[string]interface{}

func (input Dynamic) Set(value interface{}, path ...string) Dynamic {
	field, path := path[len(path)-1], path[:len(path)-1]
	current := input
	for _, segment := range path {
		next := current[segment]
		if next == nil {
			next = make(Dynamic)
			current[segment] = next
		}
		current = next.(Dynamic)
	}
	current[field] = value
	return current
}

func (input Dynamic) Get(path ...string) (interface{}, error) {
	field, rest := path[len(path)-1], path[:len(path)-1]
	current := input
	for _, segment := range rest {
		next := current[segment]
		if next == nil {
			return nil, errors.New("Path does not exist")
		}
		current = next.(Dynamic)
	}
	return current[field], nil
}

func (input Dynamic) Keys() []string {
	result := []string{}
	for key, _ := range input {
		result = append(result, key)
	}
	return result
}

func (this Dynamic) Inflate() {
	for key, value := range this {
		splits := strings.Split(key, ".")
		if casted, ok := value.(map[string]interface{}); ok {
			Dynamic(casted).Inflate()
		}
		delete(this, key)
		this.Set(value, splits...)
	}
}

func (this Dynamic) To(out interface{}) error {
	data, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return encoding.JSON.Unmarshal(bytes.NewBuffer(data), out)
}
