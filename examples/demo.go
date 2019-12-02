package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	// _ "github.com/go-sql-driver/mysql"
)

func main() {
	// 初始化数据库
	//engin, _ := gorose.Open(gorose.Config{Driver: "mysql",
	//	Dsn: "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=true"})

	// 初始化casbin适配器
	//e, _ := casbin.NewEnforcer("./model.conf", cga.NewAdapter(engin))
	e, _ := casbin.NewEnforcer("./model.conf", "./policy.csv")

	sub, obj, act := "fizz", "data1", "read"

	// Check permission.
	//e.AddPolicy(sub, obj, act)
	b,err := e.Enforce(sub, obj, act)
	fmt.Println(b,err)

	b2,err2 := e.Enforce(sub, obj, "write")
	fmt.Println(b2,err2)

	// Load the policy from DB
	err = e.LoadPolicy()
	fmt.Println(err)

	//e.AddPolicy()
	//
	//e.RemovePolicy()
	//
	//e.SavePolicy()
}
