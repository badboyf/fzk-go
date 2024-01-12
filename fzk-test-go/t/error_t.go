package main

import "fmt"

func errT1() {
	fmt.Println("run f1")
}
func errT2() {
	fmt.Println("打开数据库连接...")
	defer func(){
		err := recover()
		fmt.Println("recover: ", err)
		fmt.Println("释放数据库连接...")
	}()
	panic("出现严重错误！")
	fmt.Println("run after")
}
func errT3() {
	fmt.Println("run f3")
}
func main() {
	errT1()
	errT2()
	errT3()
}