package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Run() {
	//CasbinTest()
	//rbac_with_multiple_policy_policy()
	CasbinDB()
}

func CasbinTest() {
	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	user, err := e.GetRolesForUser("u1")
	fmt.Printf("%s\n", user)

}

func rbac_with_multiple_policy_policy() {
	e, err := casbin.NewEnforcer("./casbin/rbac_with_multiple_policy_model.conf", "./casbin/rbac_with_multiple_policy_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	fmt.Printf("%s\n", e.GetNamedPermissionsForUser("p", "user"))
	fmt.Printf("%s\n", e.GetNamedPermissionsForUser("p2", "user"))
}

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func CasbinACL() {
	//e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/policy.csv")
	e, err := casbin.NewEnforcer("/Users/fengzhikui/data/fzknotebook/fzk-go2/casbin/model.conf", "/Users/fengzhikui/data/fzknotebook/fzk-go2/casbin/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	rules := []string{"r1", "data3", "read"}
	ok, err := e.AddPolicy(rules)
	if !ok {
		fmt.Printf("%s %s", ok, err)
	}

	ok, err = e.AddRoleForUser("u1", "r1")
	if !ok {
		fmt.Printf("%s %s", ok, err)
	}
	check(e, "bbb", "data1", "read")
	check(e, "aaa", "data2", "write")
	check(e, "fzk", "data1", "write")
	check(e, "fzk", "data2", "read")
	check(e, "fzk", "data1", "read")
	check(e, "fzk", "data2", "write")
	check(e, "u1", "data3", "read")
	check(e, "r1", "data3", "read")
}

func CasbinDB() {
	a, _ := gormadapter.NewAdapter("mysql", "fzk:123456@tcp(192.168.50.157:3306)/fzk", true)
	e, _ := casbin.NewEnforcer("./casbin/model.conf", a)

	//ok, err := AddPolicy(e, "r1", "data3", "read")
	//
	//ok, err = e.AddRoleForUser("u1", "r1")
	//if !ok {
	//	fmt.Printf("%s %s", ok, err)
	//}
	//
	//check(e, "fzk", "data1", "read")
	//check(e, "u1", "data3", "read")
	//check(e, "r1", "data3", "read")
	//check(e, "u1", "data3", "write")
	//check(e, "r1", "data3", "write")

	check(e, "u2", "data1", "read")
}

func AddPolicy(e *casbin.Enforcer, sub, obj, act string) (bool, error) {
	rules := []string{sub, obj, act}
	return e.AddPolicy(rules)
}

func AddRoleForUser(e *casbin.Enforcer, u, r string) (bool, error) {
	return e.AddRoleForUser(u, r)
}
