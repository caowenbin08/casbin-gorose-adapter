package casbin_gorose_adapter

import (
	// _ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestNewAdapter(t *testing.T) {
	// 因为需要链接自己的数据库，所以这里注释掉，免得因为链接错误test不通过
	// 具体可以到 casbin-gorose-adapter/examples 目录查看用例

	// 初始化数据库
	//var engin *gorose.Engin
	//engin, _ = gorose.Open(&gorose.Config{Driver: "sqlite",
	//	Dsn: "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=true"})

	//// 初始化casbin适配器
	//e, _ := casbin.NewEnforcer("./examples/model.conf", NewAdapter(engin))
	//
	//sub, obj, act := "fizz", "data1", "read"
	//
	//// Check permission.
	//b, i := e.Enforce(sub, obj, act)
	//
	//t.Log(b,i)
	//
	//// Load the policy from DB
	//e.LoadPolicy()
}
