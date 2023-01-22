//go:build integration

package transaction

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/utilities"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestITGetTransactionByPocketId(t *testing.T) {
	//arrange
	e := echo.New()
	db := utilities.InitTestDb(t)
	_ = utilities.SeedTransactions(t, db)
	handler := New(db)

	e.GET("/cloud_pockets/:id/transactions", handler.GetTransactionByPocketId)
	expectedR := `[{"id":1,"type":"deposit","status":"success","sourcePocketId":1,"destinationPocketId":2,"description":"","amount":10,"currency":"THB","createdAt":"2021-09-01T00:00:00Z"}]`

	//action
	req := httptest.NewRequest(http.MethodGet, "/cloud_pockets/1/transactions", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	//assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedR, strings.TrimSpace(rec.Body.String()))
}
