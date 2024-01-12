package main

import "fmt"

type Person interface {
	name() string
}

type Human interface {
	Person
}

type SZ struct {
	Human
}

func (h *SZ) name() string {
	return "zhangsan"
}

func (h *SZ) sing() string {
	fmt.Printf("%v sing", h.name())

	return h.name() + " sing"
}


type User struct {
	Name   string
}
func (user *User) Login() {
	fmt.Printf("%v Login to the system", user.Name)
}

type Manager struct {
	User
	Gender string
	//Name   string
}
type Doctor struct {
	*User
}

func testUser()  {
	u := &User{Name: "foo"}
	m := Manager{
		User: *u,
		Gender: "man",
		//Name: "manager",
	}
	d := Doctor {
		User: u,
	}
	u.Name = "update"
	fmt.Printf("1  m.Name %v m.User.Name %v \n", m.Name, m.User.Name)
	fmt.Printf("2  %v %v \n", d.Name, d.User.Name)

	fmt.Println()
}

func main() {
	testUser()
}