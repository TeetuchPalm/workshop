//go:build unit

package transaction

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionByPocketId(t *testing.T) {
	//arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/:id/transactions", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	var mock sqlmock.Sqlmock
	var err error
	ti, _ := time.Parse("2021-09-01T00:00:00Z", "2001-09-28T01:00:00Z")
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
	rt := []Transaction{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM transactions WHERE sourcePocketId = $1 OR destinationPocketId = $1 ORDER BY id`))

	prep.ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "status", "sourcePocketID", "destinationPocketID", "description", "amount", "currency", "createdAt"}).AddRow(1, md.Type, md.Status, md.SourcePocketID, md.DestinationPocketID, md.Description, md.Amount, md.Currency, md.CreatedAt))

	//action
	h := New(db)
	err = h.GetTransactionByPocketId(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Greater(t, len(rt), 0)
	assert.Equal(t, 1, rt[0].ID)
	assert.Equal(t, TransactionType("deposit"), rt[0].Type)
	assert.Equal(t, TransactionStatus("success"), rt[0].Status)
	assert.Equal(t, 1, rt[0].SourcePocketID)
	assert.Equal(t, 2, rt[0].DestinationPocketID)
	assert.Equal(t, "", rt[0].Description)
	assert.Equal(t, float64(10), rt[0].Amount)
	assert.Equal(t, pocket.Currency("THB"), rt[0].Currency)
	assert.Equal(t, ti, rt[0].CreatedAt)

}
