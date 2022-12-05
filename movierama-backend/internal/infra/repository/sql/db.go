package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// DBConfig Database Configuration.
type DBConfig struct {
	Username string
	Pass     string
	Reader   string
	Writer   string
	Port     string
	DB       string
}

// NewDB returns a handle to a Database.
func NewDB(c DBConfig) (read *sql.DB, write *sql.DB) {
	read, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True",
			c.Username, c.Pass, c.Reader, c.Port, c.DB))

	if err != nil || read.Ping() != nil {
		log.Fatalf("failed to connect to master database %v", err)
		os.Exit(1)
	}

	write, err = sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True",
			c.Username, c.Pass, c.Writer, c.Port, c.DB))

	if err != nil || write.Ping() != nil {
		log.Fatalf("failed to connect to replica database %v", err)
		os.Exit(1)
	}

	return read, write
}
