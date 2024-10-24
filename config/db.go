// db.go
package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB 연결 함수
func ConnectDB() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "111111"
	dbName := "cds"
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPass, dbName))
	if err != nil {
		log.Printf("DB 연결 실패: %v", err)
		return nil, err
	}
	return db, nil
}
