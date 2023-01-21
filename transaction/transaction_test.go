//go:build unit

package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloud_pocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionById(t *testing.T) {
	//arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/transactions/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	var mock sqlmock.Sqlmock
	var err error
	ti, err := time.Parse("2021-09-01T00:00:00Z", "2001-09-28T01:00:00Z")

	if err != nil {
		fmt.Println(err)
	}
	md := Transaction{
		ID:                  1,
		Type:                "deposit",
		Status:              "success",
		SourcePocketID:      1,
		DestinationPocketID: 2,
		Description:         "",
		Amount:              10.00,
		Currency:            "THB",
		CreatedAt:           ti,
	}
	rt := Transaction{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM transactions WHERE id = $1`))

	prep.ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "status", "sourcePocketID", "destinationPocketID", "description", "amount", "currency", "createdAt"}).AddRow(1, md.Type, md.Status, md.SourcePocketID, md.DestinationPocketID, md.Description, md.Amount, md.Currency, md.CreatedAt))

	//action
	h := New(db)
	err = h.GetTransactionById(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, rt.ID)
	assert.Equal(t, TransactionType("deposit"), rt.Type)
	assert.Equal(t, TransactionStatus("success"), rt.Status)
	assert.Equal(t, 1, rt.SourcePocketID)
	assert.Equal(t, 2, rt.DestinationPocketID)
	assert.Equal(t, "", rt.Description)
	assert.Equal(t, float64(10), rt.Amount)
	assert.Equal(t, cloud_pocket.Currency("THB"), rt.Currency)
	assert.Equal(t, ti, rt.CreatedAt)

}
