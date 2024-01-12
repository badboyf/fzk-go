package main

import (
	"fmt"
)

func main() {
	pointT1()
}

func pt1()  {
	arrSlice := make([]int, 4)
	fmt.Printf("the pointer is : %p %v \n", arrSlice, arrSlice)
	fmt.Printf("the pointer is : %p %v \n", &arrSlice, arrSlice)
	printPointer(arrSlice)
	printPointer2(&arrSlice)
}

func printPointer(any []int) {
	fmt.Printf("printPointer is : %p \n", any)
	fmt.Printf("printPointer is : %p \n", &any)
}

func printPointer2(any *[]int) {
	fmt.Printf("printPointer2 is : %p \n", any)
	fmt.Printf("printPointer2 is : %p \n", &any)
}

func pointT1() {
	a1 := 1
	a2 := a1
	a3 := &a1
	a1 = 2
	fmt.Printf("a1=%v a2=%v a3=%v a3=%v \n", a1, a2, *a3, a3)
	*a3 = 3
	//a3 = 3  // 报错 int不能赋值指针
	fmt.Printf("a1=%v a2=%v a3=%v a3=%v \n", a1, a2, *a3, a3)
	changeV(a1, 5)
	fmt.Printf("changeV a1=%v a2=%v a3=%v a3=%v \n", a1, a2, *a3, a3)
	changeP(a3, 5)
	fmt.Printf("changeP a1=%v a2=%v a3=%v a3=%v \n", a1, a2, *a3, a3)

	s1 := [10]int{1,2,3,4,5,6,7,8,9,0}
	s3 := s1
	s5 := &s1
	s2 := []int{1,2,3,4,5,6,7,8,9,0}
	s4 := s2
	s6 := &s2
	s1[0] = 6
	s5[1] = 8
	(*s5)[2] = 8
	s2[0] = 6
	(*s6)[1] = 8
	//s6[1] = 8
	fmt.Printf("类型 \t s1=%T s3=%T s5=%T \n\t s2=%T s4=%T s6=%T \n", s1, s3, s5, s2, s4, s6)
	fmt.Printf("init \t s1=%v s3=%v s5=%v \n\t s2=%v s4=%v s6=%v \n", s1, s3, s5, s2, s4, s6)
	changeArray(s1, 0)
	changeSlice(s2, 0)
	fmt.Printf("change \t s1=%v s3=%v s5=%v \n\t s2=%v s4=%v s6=%v \n", s1, s3, s5, s2, s4, s6)
	changeSliceP(&s2, 0)
	fmt.Printf("sliceP \t s1=%v s3=%v s5=%v \n\t s2=%v s4=%v s6=%v \n", s1, s3, s5, s2, s4, s6)
}

func changeV(a1 int, a3 int) {
	a1 = a3
}
func changeP(a2 *int, a3 int) {
	*a2 = a3
}

func changeSlice(s []int, v int)  {
	s[4] = v
}
func changeArray(s [10]int, v int)  {
	s[4] = v
}
func changeSliceP(s *[]int, v int)  {
	(*s)[5] = v
}