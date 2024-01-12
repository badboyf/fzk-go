package util

import "strings"

var DB = &dbUtil{}
var DBStr = &dbStringUtil{}

type dbUtil struct {
	path    string
	Content []*map[string]interface{}
}
type dbStringUtil struct {
	path    string
	Content []string
}

func (a dbUtil) Init(filename string) dbUtil {
	content := File.ReadFilePanic(filename)
	arr := Json.ToArray([]byte(content))

	return dbUtil{
		path:    filename,
		Content: arr,
	}
}

func (a dbUtil) StoreAndFlushOne(item interface{}) error {
	v := Json.ToMapNotNull(Json.ToByte(item))
	a.Content = append(a.Content, v)
	content := Json.ToStringPretty(a.Content)
	err := File.Write(a.path, content)
	return err
}


func (a dbStringUtil) Init(filename string) dbStringUtil {
	content := File.ReadFilePanic(filename)
	arr := strings.Split(content, "\n")

	return dbStringUtil{
		path:    filename,
		Content: arr,
	}
}

func (a dbStringUtil) StoreAndFlushOne(item string) error {
	err := File.WriteAppend(a.path, item + "\n")
	return err
}
