//go:build integration
// +build integration

package sql

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDB_CreateConnection(t *testing.T) {
	godotenv.Load("../../../../.env.dist")

	c := DBConfig{
		MySQLUsername: getEnv("MYSQL_USERNAME"),
		MySQLPass:     getEnv("MYSQL_PASSWORD"),
		MySQLReader:   getEnv("MYSQL_READ"),
		MySQLWriter:   getEnv("MYSQL_WRITE"),
		MySQLPort:     getEnv("MYSQL_PORT"),
		MySQLDB:       getEnv("MYSQL_DB"),
	}

	r, w := NewDB(c)
	assert.IsType(t, &sql.DB{}, r)
	assert.IsType(t, &sql.DB{}, w)
	assert.NoError(t, r.Ping())
	assert.NoError(t, w.Ping())

}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Missing configuration %s", key)
	}
	return v
}
