package session

import (
	"database/sql"
	"geeorm/log"
	"strings"
)

// Session 保持一个到 sql.DB 的指针，提供所有的数据库操作
type Session struct {
	db      *sql.DB
	sql     strings.Builder // sql 语句
	sqlVars []interface{}   // sql 语句中对应的占位符
}

// Session 的构造函数
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// 清除 session 的状态
// 将 sql 语句和，参数（占位符）清空
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// DB 返回 *sql.DB
func (s *Session) DB() *sql.DB {
	return s.db
}

// 执行 sql 脚本
// 执行完后清除 session 状态
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 从 db 获取一条记录
// 执行完后清除 session 状态
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 从 db 获取记录集合
// 执行完后清除 session 状态
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// Appends sql and sqlVars
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}
