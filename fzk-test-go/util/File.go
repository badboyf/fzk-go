package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var File = &fileUtil{}

type fileUtil struct {
}

func (a fileUtil) ReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (a fileUtil) ReadFilePanic(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (a fileUtil) Write(filename string, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0666)
	return err
}

func (a fileUtil) WriteAppend(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}

	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(content)

	write.Flush()
	return err
}