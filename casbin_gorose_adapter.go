package casbin_gorose_adapter

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gohouse/gorose/v2"
)

// Adapter is the file adapter for Casbin.
// It can load policy from file or save policy to file.
type CasbinGoroseAdapter struct {
	*gorose.Engin
}

// CasbinRule casbin_rule table
type CasbinRule struct {
	Id       	string `gorose:"id"`
	PType       string `gorose:"p_type"`
	V0          string `gorose:"v0"`
	V1          string `gorose:"v1"`
	V2          string `gorose:"v2"`
	V3          string `gorose:"v3"`
	V4          string `gorose:"v4"`
	V5          string `gorose:"v5"`
}
// TableName return table's true name
func (*CasbinRule) TableName() string {
	return "casbin_rule"
}

// NewCasbinGoroseAdapter is the constructor for Adapter.
func NewAdapter(ge *gorose.Engin) *CasbinGoroseAdapter {
	cga := &CasbinGoroseAdapter{ge}
	// 建表，如果是mysql驱动的话
	if err := cga.createTable(); err != nil {
		panic(err.Error())
	}
	return cga
}

func (a *CasbinGoroseAdapter) createTable() (err error) {
	// 如果传入的驱动是mysql，则执行此操作
	if a.Engin.GetDriver() == gorose.DriverMysql {
		sqlStr := `CREATE TABLE IF NOT EXISTS casbin_rule (
  id int(11) NOT NULL AUTO_INCREMENT,
  p_type varchar(32) NOT NULL DEFAULT '' COMMENT 'perm类型：p,g......',
  v0 varchar(64) NOT NULL DEFAULT '' COMMENT '角色名字...',
  v1 varchar(64) NOT NULL DEFAULT '' COMMENT '对象资源...',
  v2 varchar(64) NOT NULL DEFAULT '' COMMENT '权限值...',
  v3 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
  v4 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
  v5 varchar(64) NOT NULL DEFAULT '' COMMENT 'ext',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
		_, err = a.Engin.NewSession().Execute(sqlStr)
	}

	return
}

// LoadPolicy loads all policy rules from the storage.
func (a *CasbinGoroseAdapter) LoadPolicy(model model.Model) error {
	if a.Engin == nil {
		return errors.New("invalid gorose.Engin, gorose.Engin cannot be empty")
	}

	return a.loadPolicyData(model, persist.LoadPolicyLine)
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinGoroseAdapter) SavePolicy(model model.Model) error {
	if a.Engin == nil {
		return errors.New("db error")
	}

	var cr []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			var crrow = a.builCasbinRule(ptype,rule)
			cr = append(cr, crrow)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			var crrow = a.builCasbinRule(ptype,rule)
			cr = append(cr, crrow)
		}
	}

	aff, e := a.Engin.NewOrm().Insert(cr)
	if e!=nil {
		return e
	}
	if aff==0 {
		return errors.New("insert db fail")
	}
	return nil
}

func (a *CasbinGoroseAdapter) loadPolicyData(model model.Model, handler func(string, model.Model)) error {

	// get rules
	var cr []CasbinRule
	if e := a.Engin.NewOrm().Table(&cr).Select(); e!=nil {
		return e
	}

	for _, item := range cr {
		var line = item.PType
		if item.V0 != "" { line = fmt.Sprintf("%s, %s", line, item.V0) }
		if item.V1 != "" { line = fmt.Sprintf("%s, %s", line, item.V1) }
		if item.V2 != "" { line = fmt.Sprintf("%s, %s", line, item.V2) }
		if item.V3 != "" { line = fmt.Sprintf("%s, %s", line, item.V3) }
		if item.V4 != "" { line = fmt.Sprintf("%s, %s", line, item.V4) }
		if item.V5 != "" { line = fmt.Sprintf("%s, %s", line, item.V5) }
		handler(line, model)
	}
	return nil
}

func (a *CasbinGoroseAdapter) builCasbinRule(ptype string, rule []string) (cr CasbinRule) {
	var length = len(rule)
	if length == 0 {
		return
	}
	cr.PType = ptype
	switch {
	case length>0: cr.V0 = rule[0]; fallthrough
	case length>1: cr.V1 = rule[1]; fallthrough
	case length>2: cr.V2 = rule[2]; fallthrough
	case length>3: cr.V3 = rule[3]; fallthrough
	case length>4: cr.V4 = rule[4]; fallthrough
	case length>5: cr.V5 = rule[5]
	}
	return
}

// AddPolicy adds a policy rule to the storage.
func (a *CasbinGoroseAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	var cr = a.builCasbinRule(ptype, rule)
	// insert to db
	aff,err := a.Engin.NewOrm().Insert(&cr)
	if err!=nil {
		return err
	}
	if aff==0 {
		return errors.New("insert db error")
	}
	return nil
}

// RemovePolicy removes a policy rule from the storage.
func (a *CasbinGoroseAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	var cr = a.builCasbinRule(ptype, rule)
	return a.rawDelete(&cr)
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *CasbinGoroseAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {

	var line CasbinRule

	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	return a.rawDelete(&line)
}

func (a *CasbinGoroseAdapter) rawDelete(line *CasbinRule) error {
	var where = make(map[string]string)
	where["p_type"] = line.PType

	if line.V0 != "" { where["v0"] = line.V0 }
	if line.V1 != "" { where["v1"] = line.V1 }
	if line.V2 != "" { where["v2"] = line.V2 }
	if line.V3 != "" { where["v3"] = line.V3 }
	if line.V4 != "" { where["v4"] = line.V4 }
	if line.V5 != "" { where["v5"] = line.V5 }

	aff,err := a.Engin.NewOrm().Where(where).Delete()
	if err!=nil {
		return err
	}
	if aff==0 {
		return errors.New("delete fail")
	}
	return nil
}