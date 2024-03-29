package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type jsonContent struct {
	data      map[string]interface{}
	sliceData []interface{}
}

func (c *jsonContent) GetBody() ([]byte, error) {
	if len(c.data) > 0 {
		return json.Marshal(c.data)
	} else if len(c.sliceData) > 0 {
		return json.Marshal(c.sliceData)
	} else {
		return nil, nil
	}
}

func (c *jsonContent) SetBody(content interface{}) error {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		return c.setJsonStringBody(v.String())
	case reflect.Struct:
		return c.setStructBody(v.Interface())
	case reflect.Slice:
		c.setSliceBody(v)
	default:
		return errors.New("content is not valid json")
	}
	return nil
}

func (c *jsonContent) setJsonStringBody(content string) error {
	var val interface{}
	d := json.NewDecoder(strings.NewReader(content))
	d.UseNumber()
	if err := d.Decode(&val); err == nil {
		switch v := reflect.ValueOf(val); v.Kind() {
		case reflect.Map:
			return c.setStructBody(val)
		case reflect.Slice:
			c.setSliceBody(v)
		default:
			return errors.New("conent is not valid json")
		}
		return nil
	} else {
		return err
	}
}

func (c *jsonContent) setStructBody(content interface{}) error {
	if marshalContent, err := json.Marshal(content); err != nil {
		return err
	} else {
		c.data = make(map[string]interface{})
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&c.data); err != nil {
			return err
		}
	}
	return nil
}

func (c *jsonContent) setSliceBody(v reflect.Value) {
	c.sliceData = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		c.sliceData[i] = v.Index(i).Interface()
	}
}
