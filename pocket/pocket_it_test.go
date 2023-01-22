//go:build integration

package pocket

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/utilities"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestITGetOne(t *testing.T) {
	//arrangement
	e := echo.New()

	db := utilities.InitTestDb(t)
	id := utilities.SeedPocket(t, db)
	handler := New(db)

	e.GET("/cloud-pockets/:id", handler.GetOne)
	expectedR := `{"id":1,"name":"demoPocket","category":"test","amount":100,"goal":20000.02,"currency":"THB","createdAt":"2021-09-01T00:00:00Z","updatedAt":"2021-09-01T00:00:00Z"}`

	//action
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets/"+strconv.Itoa(int(id)), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	t.Log(rec.Body.String())

	//assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedR, strings.TrimSpace(rec.Body.String()))

}
