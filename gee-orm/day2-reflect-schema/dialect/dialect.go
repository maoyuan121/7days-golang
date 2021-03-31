package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// Dialect is an interface contains methods that a dialect has to implement
type Dialect interface {
	// 根据 struct 的字段类型得出数据库得数据类型
	DataTypeOf(typ reflect.Value) string

	// 判断数据库中是否存在指定的表
	// 返回值的 string 是 sql 脚本
	TableExistSQL(tableName string) (string, []interface{})
}

// 注册一个 dialect 到全局变量中
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 从全局变量中获取获取指定的 dialect
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
