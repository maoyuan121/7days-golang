package clause

import (
	"strings"
)

// SQL 子句
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// 子句类型
type Type int

// 支持的子句类型
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// 添加一个特殊类型的子句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// 生成最终的 SQL 和 SQL 参数
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
