# casbin-gorose-adapter
the gorose adapter for casbin

casbin-gorose-adapter is the gorose adapter for Casbin. With this library, Casbin can load policy from gorose supported database or save policy to it.

Based on Officially Supported Databases, The current supported databases are:

MySQL  
PostgreSQL  
Sqlite3  
SQL Server  

## Installation
```shell script
go get github.com/gohouse/casbin-gorose-adapter
```

## simple example
```go
package main

import (
	"github.com/casbin/casbin/v2"
	cga "github.com/gohouse/casbin-gorose-adapter"
	"github.com/gohouse/gorose/v2"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// init db orm
	engin, _ := gorose.Open(&gorose.Config{Driver: "mysql", 
		Dsn: "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=true"})

	// init casbin adapter
	// model can use in source code as `github.com/casbin/casbin/v2/model.NewModel().AddDef("r","r","sub, obj, act")`
	e, _ := casbin.NewEnforcer("./model.conf", cga.NewAdapter(engin))

	sub, obj, act := "fizz", "data1", "read"

	// Check permission.
	e.Enforce(sub, obj, act)

	// Load the policy from DB
	e.LoadPolicy()

	//e.AddPolicy()
	//
	//e.RemovePolicy()
	//
	//e.SavePolicy()
}
```

## generate table
now, only mysql driver can auto generate table, the table sql like :  

```sql
CREATE TABLE IF NOT EXISTS casbin_rule (
 id int(11) NOT NULL AUTO_INCREMENT,
 p_type varchar(32) NOT NULL DEFAULT '' COMMENT 'perm类型：p,g......',
 v0 varchar(64) NOT NULL DEFAULT '' COMMENT '角色名字(rbac),其他角色类型则依次存放v0-v5...',
 v1 varchar(64) NOT NULL DEFAULT '' COMMENT '对象资源(rbac),其他角色类型则依次存放v0-v5...',
 v2 varchar(64) NOT NULL DEFAULT '' COMMENT '权限值(rbac),其他角色类型则依次存放v0-v5...',
 v3 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
 v4 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
 v5 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
 PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='casbin权限规则表';
```