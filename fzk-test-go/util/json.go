package util

import (
	"bytes"
	"encoding/json"
)

var Json = &jsonInner{}

type jsonInner struct {
}

// ToStringPretty 转成string
func (t *jsonInner) ToStringPretty(i interface{}) string {
	data := t.ToString(i)
	var out bytes.Buffer
	err := json.Indent(&out, []byte(data), "", "  ")
	if err != nil {
		return data
	}
	return out.String()
}

// ToString 转成string
func (t *jsonInner) ToString(i interface{}) string {
	return string(t.ToByte(i))
}

// ToByte 转成byte数组
func (t *jsonInner) ToByte(i interface{}) []byte {
	rs, err2 := json.Marshal(i)
	if err2 != nil {
		return []byte("")
	}
	return rs
}

// ToBean 转为成bean
func (t *jsonInner) ToBean(src interface{}, i interface{}) (e error) {
	if src == nil {
		return nil
	}

	if _, ok := src.(string); ok {
		e = json.Unmarshal([]byte(src.(string)), &i)
	} else if _, ok := src.([]byte); ok {
		e = json.Unmarshal(src.([]byte), &i)
	} else {
		marshal, e := json.Marshal(src)
		if e != nil {
			return e
		}
		e = json.Unmarshal(marshal, &i)
	}

	return
}

// ToMapNotNull 转为map
func (t *jsonInner) ToMapNotNull(v []byte) *map[string]interface{} {
	var m *map[string]interface{}
	if v == nil {
		return m
	}
	_ = json.Unmarshal(v, &m)
	return m
}

// ToArray 转为数组
func (t *jsonInner) ToArray(v []byte) []*map[string]interface{} {
	var m []*map[string]interface{}
	if v == nil {
		return m
	}
	_ = json.Unmarshal(v, &m)
	return m
}

func (t *jsonInner) Merge(v1 *map[string]interface{}, v2 *map[string]interface{}) *map[string]interface{} {
	if v2 == nil {
		return v1
	} else if v1 == nil {
		return v2
	} else {
		for k := range *v2 {
			(*v1)[k] = (*v2)[k]
		}
		return v1
	}
}

