package db

import (
	"database/sql"
	"fmt"

	"github.com/ability-sh/abi-micro/micro"
)

type DBExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type DBInterface interface {
	DBExec
	Begin() (DBTransaction, error)
}

type DBTransaction interface {
	DBExec
	Commit() error
	Rollback() error
}

type DBService interface {
	micro.Service
	GetDB() *sql.DB
}

func GetDB(ctx micro.Context, name string) (*sql.DB, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(DBService)
	if ok {
		ctx.AddCount("db", 1)
		return ss.GetDB(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof DBService", name)
}

func GetDBInterface(ctx micro.Context, name string) (DBInterface, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(DBService)
	if ok {
		ctx.AddCount("db", 1)
		return &dbInterface{db: ss.GetDB(), ctx: ctx}, nil
	}
	return nil, fmt.Errorf("service %s not instanceof DBService", name)
}

type dbInterface struct {
	ctx micro.Context
	db  *sql.DB
}

func (db *dbInterface) Exec(query string, args ...interface{}) (sql.Result, error) {
	st := db.ctx.Step("db.Exec")
	rs, err := db.db.Exec(query, args...)
	if err != nil {
		st("[err:1] [sql:%s] %s", query, err.Error())
	} else {
		st("")
	}
	return rs, err
}

func (db *dbInterface) Query(query string, args ...interface{}) (*sql.Rows, error) {
	st := db.ctx.Step("db.Query")
	rs, err := db.db.Query(query, args...)
	if err != nil {
		st("[err:1] [sql:%s] %s", query, err.Error())
	} else {
		st("")
	}
	return rs, err
}

func (db *dbInterface) Begin() (DBTransaction, error) {
	st := db.ctx.Step("db.Begin")
	tx, err := db.db.Begin()
	if err != nil {
		st("[err:1] %s", err.Error())
	} else {
		st("")
	}
	return &dbTransaction{tx: tx, ctx: db.ctx}, err
}

type dbTransaction struct {
	tx  *sql.Tx
	ctx micro.Context
}

func (db *dbTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	st := db.ctx.Step("db.tx.Exec")
	rs, err := db.tx.Exec(query, args...)
	if err != nil {
		st("[err:1] [sql:%s] %s", query, err.Error())
	} else {
		st("")
	}
	return rs, err
}

func (db *dbTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	st := db.ctx.Step("db.tx.Query")
	rs, err := db.tx.Query(query, args...)
	if err != nil {
		st("[err:1] [sql:%s] %s", query, err.Error())
	} else {
		st("")
	}
	return rs, err
}

func (db *dbTransaction) Commit() error {
	st := db.ctx.Step("db.tx.Commit")
	err := db.tx.Commit()
	if err != nil {
		st("[err:1] %s", err.Error())
	} else {
		st("")
	}
	return err
}

func (db *dbTransaction) Rollback() error {
	st := db.ctx.Step("db.tx.Rollback")
	err := db.tx.Rollback()
	if err != nil {
		st("[err:1] %s", err.Error())
	} else {
		st("")
	}
	return err
}
