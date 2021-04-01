package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// 代表数据库表中的一个列
type Field struct {
	Name string // 列名
	Type string // 数据类型
	Tag  string // 对应的  tag
}

// 代表数据库中的一个表
type Schema struct {
	Model      interface{}       // 对应的 Model struct
	Name       string            // 表名
	Fields     []*Field          // 列集合
	FieldNames []string          // 列名集合
	fieldMap   map[string]*Field // 列名和列的映射
}

// 根据列名返回列
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 返回 dest 的成员变量的值
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

type ITableName interface {
	TableName() string
}

// 将一个 struct 解析到对应的 dialect 的 schema
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 如果 dest 实现了 ITableName 那么表名是其 TableName() 函数的返回值，否在就是其 struct name
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)

		// 如果字段不是匿名的且是导出（public，大写开头）的
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
