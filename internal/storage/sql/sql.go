package sql

import (
	"database/sql"
	"student-api/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type Sql struct {
	Db *sql.DB
}

func New(cfg config.Config) (*Sql, error) {
	db, err := sql.Open("mysql", cfg.Storage_path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY, 
		name TEXT, 
		email TEXT, 
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sql{Db: db}, nil
}
