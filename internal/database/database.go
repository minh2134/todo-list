package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"todo-list/internal/task"

	_ "github.com/mattn/go-sqlite3"
)

// Database struct to hide sqlite3 Go complexity
// consists of a read-only conn pool, transactions are configured with BEGIN DEFERRED (default)
// and a read-write conn pool, max conns is locked to 1 (all serialized)
// and with a BEGIN IMMEDIATE transaction config
// source: https://github.com/mattn/go-sqlite3/issues/1022
type Database struct {
	readOnly  *sql.DB
	readWrite *sql.DB
}

var (
	ErrorAddingTransaction error = errors.New("Something wrong happened with the transaction")
)

func Open(filename string) (Database, error) {
	var (
		db  Database
		err error
	)

	if filename == "" {
		err = errors.New("Empty filename")
		return db, err
	}

	// database general config
	dbOptsURL := "file:" + filename + "?_journal=WAL"
	dbType := "sqlite3"
	maxConnLifetime := 5 * time.Minute
	maxIdleConnLifetime := 5 * time.Minute

	// read-only conn pool specific config
	maxOpenConnRO := 10
	maxIdleConnRO := 5

	// read-write conn pool specific config
	dbOptsURLRW := dbOptsURL + "&_txlock=immediate&_timeout=5000"
	maxOpenConnRW := 1
	maxIdleConnRW := 1

	// read-only pool
	db.readOnly, err = sql.Open(dbType, dbOptsURL)
	if err != nil {
		return db, err
	}
	db.readOnly.SetMaxIdleConns(maxIdleConnRO)
	db.readOnly.SetMaxOpenConns(maxOpenConnRO)
	db.readOnly.SetConnMaxIdleTime(maxIdleConnLifetime)
	db.readOnly.SetConnMaxLifetime(maxConnLifetime)

	// read-write pool
	db.readWrite, err = sql.Open(dbType, dbOptsURLRW)
	if err != nil {
		return db, err
	}
	db.readWrite.SetMaxIdleConns(maxIdleConnRW)
	db.readWrite.SetMaxOpenConns(maxOpenConnRW)
	db.readWrite.SetConnMaxIdleTime(maxIdleConnLifetime)
	db.readWrite.SetConnMaxLifetime(maxConnLifetime)

	// Always run bootstrap so new tables are created if they don't exist yet
	// (uses CREATE TABLE IF NOT EXISTS — safe to run on existing databases)
	db.bootstrap()

	return db, err
}

func (db Database) Close() {
	db.readOnly.Close()
	db.readWrite.Close()
}

func (db Database) bootstrap() error {
	conn := db.readWrite

	schema := `
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	description TEXT,
	completed BOOLEAN
);
	`

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(schema)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (db Database) InsertTask(t task.Task) (int64, error) {
	conn := db.readWrite

	queryTransaction := `
INSERT INTO tasks (name, description, completed) VALUES (?,?,?);
	`
	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		queryTransaction,
		t.Name,
		t.Desc,
		t.Completed,
	)
	if err != nil {
		return 0, err
	}

	tId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return tId, err
}

func (db Database) GetTasks(lq ListQuery) (map[int]task.Task, error) {
	var (
		query strings.Builder
		args  []any
	)
	query.WriteString(`SELECT * FROM tasks WHERE 1=1`)
	if lq.Name != "" {
		query.WriteString(" AND name LIKE ?")
		// search all matching substrings instead of exact
		args = append(args, "%"+lq.Name+"%")
	}
	if lq.Completed != ALL {
		query.WriteString(" AND completed=?")
		args = append(args, lq.Completed)
	}
	fmt.Println(query.String(), args)

	tsks := make(map[int]task.Task)

	conn := db.readOnly

	tx, err := conn.Begin()
	if err != nil {
		return tsks, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(query.String(), args...)
	if err != nil {
		return tsks, err
	}

	for rows.Next() {
		var (
			name      string
			desc      string
			completed bool
			id        int
		)
		err := rows.Scan(&id, &name, &desc, &completed)
		if err != nil {
			return tsks, err
		}
		tsk := task.MakeTask(name, desc)
		tsk.Completed = completed

		tsks[id] = tsk
	}

	tx.Commit()
	return tsks, err
}
