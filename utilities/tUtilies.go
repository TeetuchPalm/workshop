package utilities

import (
	"database/sql"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	_ "github.com/lib/pq"
)

func SeedTransactions(t *testing.T, db *sql.DB) int64 {
	stmt, err := db.Prepare(`INSERT INTO transactions (type, status, sourcePocketId, destinationPocketID, description, amount, currency, createdAt) VALUES ('deposit', 'success', 1, 2, '', 10.00, 'THB', '2021-09-01T00:00:00Z') RETURNING id`)
	if err != nil {
		t.Fatal(err)
	}
	row := stmt.QueryRow()
	var id int64
	row.Scan(&id)
	return id
}

func InitTestDb(t *testing.T) *sql.DB {
	cfg := config.New().All()
	db, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
