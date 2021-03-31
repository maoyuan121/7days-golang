package geeorm

import (
	"database/sql"

	"geeorm/log"
	"geeorm/session"
)

// Engine 是 geeorm 的主 strcut，管理所有的 db session 和事务
type Engine struct {
	db *sql.DB
}

// Engine 的构造函数
// 链接数据库并 ping 下看是否是通的
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

// 关闭数据库连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// 为下一个操作创建一个 session
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}
