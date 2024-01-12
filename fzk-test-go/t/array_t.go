package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func arrT4()  {
	var a1 = []string{"a", "b", "c", "d"}
	var a2 = make([]string, len(a1))
	copy(a2, a1)
	m := map[string]int{}
	for i, v := range a2 {
		m[v] = i
	}
	remove := []string{"b", "c"}
	for _, v := range remove {
		delete(m, v)
	}
	fmt.Printf("%v\n", a2)
	fmt.Printf("%v\n", m)
}

// 数组实例化方式
func arrT3()  {
	var array1 = [5]int{1, 2, 3}
	var array2 = [...]int{6, 7, 8}
	var array3 = []int{9, 10, 11, 12}
	var array4 = [5]string{3: "Chris", 4: "Ron"}
	var array5 = [...]string{3: "Tom", 2: "Alice"}
	var array6 = []string{4: "Smith", 2: "Alice"}

	fmt.Printf("%v", array1[:])
	fmt.Printf("%v", array1[0:])
	fmt.Printf("%v", array1[:2])

	slice := make([]int, 0, 100)

	fmt.Printf("%v %v %v %v %v", array2, array3, array4, array5, array6, slice)
	for i, v := range array1 {
		fmt.Printf("index:%d  value:%d\n", i, v)
	}
}

func arrT2()  {
	//a1 := []string{}
	var a1 []string
	var i interface{}
	i = a1
	fmt.Printf("%v %v  \n", i, &i)
	fmt.Println(IsNil(i))
}

// 判断 interface 空
func IsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

func arrT1()  {
	a1 := [5]int {0, 1, 2, 3, 4}
	a2 := &[...] int {1,2,3,4,5}
	a1[0] = 6
	a2[0] = 6
	fmt.Printf("%v  %v  %p \n", a1, &a1, &a1)
	fmt.Println(a2)

	for i, v := range a2 {
		fmt.Println("a2 -> ", i, "=", v)
	}
}



func ConvertToMapArray() {
	articleStrings := `[{"name":"name1","amount":100},{"name":"name2","amount":200},{"name":"name3","amount":300},{"name":"name4","amount":400}]`
	var articleSlide []map[string]interface{}
	multiErr := json.Unmarshal([]byte(articleStrings), &articleSlide)
	if multiErr != nil {
		fmt.Println("转换出错：", multiErr)
	}
	fmt.Printf("字符串转map数组 %v %v \n", articleSlide, reflect.TypeOf(articleSlide))
}

type robot struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
func ConvertToStructArray() {
	str := `[{"name":"name1","amount":100},{"name":"name2","amount":200},{"name":"name3","amount":300},{"name":"name4","amount":400}]`
	var all []robot
	err := json.Unmarshal([]byte(str), &all)
	if err != nil {
		fmt.Printf("err=%v", err)
	}
	fmt.Printf("字符串转结构体 %v %v \n", all, reflect.TypeOf(all))
}
func ConvertToStructArrayFromMapArray() {
	str := `[{"name":"name1","amount":100},{"name":"name2","amount":200},{"name":"name3","amount":300},{"name":"name4","amount":400}]`
	var mapSlide []map[string]interface{}
	multiErr := json.Unmarshal([]byte(str), &mapSlide)
	if multiErr != nil {
		fmt.Println("转换出错：", multiErr)
	}

	var rs []robot
	marshal, _ := json.Marshal(&mapSlide)
	multiErr = json.Unmarshal(marshal, &rs)
	if multiErr != nil {
		fmt.Println("转换出错：", multiErr)
	}
	fmt.Printf("map数组转结构体 %v %v %p %v \n", rs, &rs, rs, reflect.TypeOf(rs))

	var rsi interface{} = make([]*robot, 0, 1)
	multiErr = json.Unmarshal(marshal, &rsi)
	if multiErr != nil {
		fmt.Println("转换出错：", multiErr)
	}
	fmt.Printf("map数组转结构体 %v %v %p %v \n", rsi, &rsi, rsi, reflect.TypeOf(rsi))

	fmt.Println("接口转数组并修改值==========================")
	i := rsi.([]interface{})
	for _, v := range i {
		fmt.Printf("%v %T %v \n", v, v, reflect.TypeOf(v).Kind() == reflect.Map)
		m := v.(map[string]interface{})
		fmt.Printf("%v %T\n", m["amount"], m)
	}
}

func arrT5() {
	a1 := []ArrT{{Name: "1"}, {Name: "2"}}
	fmt.Printf("%v\n", a1)
	a1[0] = ArrT{Name: "1-update"}
	fmt.Printf("修改整个对象 %v\n", a1)
	a1[0].Name = "1-update-again"
	fmt.Printf("修改对象的值 %v\n", a1)

	for i := range a1 {
		a1[i].Name = strconv.Itoa(i)
	}
	fmt.Printf("循环修改对象值 %v\n", a1)

	a1p := &a1
	(*a1p)[0] = ArrT{Name: "1-update-p"}
	fmt.Printf("指针 %v  %v\n", a1p, a1)
	for i := range *a1p {
		t := (*a1p)[i]
		t.Name = fmt.Sprintf("%v %v", i, "point")
	}
	fmt.Printf("指针 %v  %v\n", a1p, a1)

	for i := range *a1p {
		t := &(*a1p)[i]
		t.Name = fmt.Sprintf("%v %v", i, "point")
	}
	fmt.Printf("指针 %v  %v\n", a1p, a1)

	a2 := *a1p
	for i := range a2 {
		a2[i].Name = fmt.Sprintf("%v%v", "a2", i)
	}
	fmt.Printf("a2 %v  %v %v\n", a1p, a1, a2)
}

type ArrT struct {
	Name string
}

func main() {
	//arrT1()
	//arrT3()
	//ConvertToMapArray()
	//ConvertToStructArray()
	//ConvertToStructArrayFromMapArray()
	//arrT4()
	arrT5()
}