package utilities

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
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

func SeedPocket(t *testing.T, db *sql.DB) int64 {
	stmt, err := db.Prepare(`INSERT INTO pockets (name, category, amount, goal, currency, createdAt, updatedAt) VALUES ('demoPocket', 'test', 100, 20000.02, 'THB', '2021-09-01T00:00:00Z', '2021-09-01T00:00:00Z') RETURNING id`)
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

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func Request(method, uri string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, uri, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
