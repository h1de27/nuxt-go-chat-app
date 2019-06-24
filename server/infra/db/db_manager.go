package db

import (
	"context"
	"database/sql"

	"github.com/hideUW/nuxt-go-chat-app/server/domain/model"
	"github.com/hideUW/nuxt-go-chat-app/server/domain/repository"

	// SQL Driver
	_ "github.com/go-sql-driver/mysql"
)

// dbManager manages SQL.
type dbManager struct {
	Conn *sql.DB
}

// NewDBManager generates and returns SQLManager.
func NewDBManager() repository.DBManager {
	conn, err := sql.Open("mysql", "root@tcp(nvgdb:3306)/nuxt-go-chat-app?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err.Error())
	}

	return &dbManager{
		Conn: conn,
	}
}

// Exec executes SQL.
func (s dbManager) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.Exec(query, args...)
}

// ExecContext executes SQL with context.
func (s *dbManager) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.ExecContext(ctx, query, args...)
}

// Query executes query which returns row.
func (s *dbManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		err = &model.SQLError{
			BaseErr:                   err,
			InvalidReasonForDeveloper: "failed to execute query",
		}
		return nil, err
	}

	return rows, nil
}

// Query executes query which returns row with context.
func (s *dbManager) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		err = &model.SQLError{
			BaseErr:                   err,
			InvalidReasonForDeveloper: "failed to execute query with context",
		}
		return nil, err
	}

	return rows, nil
}

// Prepare prepares statement for Query and Exec later.
func (s *dbManager) Prepare(query string) (*sql.Stmt, error) {
	return s.Conn.Prepare(query)
}

// Prepare prepares statement for Query and Exec later with context.
func (s *dbManager) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.Conn.PrepareContext(ctx, query)
}

// Begin begins tx.
func (s *dbManager) Begin() (repository.TxManager, error) {
	return s.Conn.Begin()
}
