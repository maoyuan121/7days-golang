package clause

import (
	"fmt"
	"strings"
)

// 生成器
type generator func(values ...interface{}) (string, []interface{})

// 生成器映射
// key 是 insert、values、select、limit、where、orderby
// 值是具体的生成器
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// 插入生成器
// 第一个参数为表名，后面的是要插入的列
func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// Values 生成器
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars

}

// Select 生成器
// 第一个参数为表名，后面的为要查询的字段
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// Limit 生成器
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

// Where 生成器
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// OrderBy 生成器
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// Update  生成器
// 总共两个参数
// 第一个参数是表名, 第二个参数是要更新的字段的键值映射
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// Delete 生成器
// 总共一个参数，即表名
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// Count 生成器
// 总共一个参数，即表名
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
