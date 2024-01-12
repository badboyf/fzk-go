package main

import (
	"fmt"
)

func chFuncWirte(num int, c chan int) {
	c <- num + 1 // 写入数据
}
func chFuncRead(c chan int) {
	// 读取数据，阻塞
	ret := <-c
	fmt.Println(ret) // 打印：2
}
func chFunc2()  {
	// 定义
	c := make(chan int)

	go chFuncWirte(1,c)
	// 读取数据，阻塞
	ret := <-c

	fmt.Println(ret) // 打印：2
}

func chFunc3()  {
	// 定义
	c := make(chan int, 2)

	c <- 1
	fmt.Printf("%v\n", <- c)

	c <- 2
	c <- 3
	fmt.Printf("%v\n", <- c)
	fmt.Printf("%v\n", <- c)

	c <- 4
	c <- 5
	c <- 6
	fmt.Printf("%v\n", <- c)
	fmt.Printf("%v\n", <- c)
	fmt.Printf("%v\n", <- c)
}


func main() {
	//chFunc2()
	chFunc3()
}