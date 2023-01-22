//go:build unit

package pocket

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetOne(t *testing.T) {
	testcases := []struct {
		name       string
		id         string
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{
			name: "Get cloud pocket successfully",
			id:   "1",
			sqlFn: func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}

				row := sqlmock.NewRows([]string{
					"id",
					"name",
					"category",
					"amount",
					"goal",
					"currency",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					1,
					"travel",
					"travel",
					1000.0,
					1000.0,
					"THB",
					"2021-09-01T00:00:00Z",
					"2021-09-01T00:00:00Z",
					nil,
				)
				mock.ExpectQuery(cGetOneStmt).WithArgs("1").WillReturnRows(row)
				return db, err
			},
			wantStatus: http.StatusOK,
			wantBody: `{
				"id": 1,
				"name": "travel",
				"category": "travel",
				"amount": 1000.0,
				"goal": 1000.0,
				"currency": "THB",
				"createdAt": "2021-09-01T00:00:00Z",
				"updatedAt": "2021-09-01T00:00:00Z"
			}`,
		},
	}

	for _, tc := range testcases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames(cID)
		c.SetParamValues(tc.id)

		db, err := tc.sqlFn()
		h := New(db)

		assert.NoError(t, err)
		if assert.NoError(t, h.GetOne(c)) {
			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		}
	}
}

func TestGetOne_Error(t *testing.T) {
	testcases := []struct {
		name       string
		id         string
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
	}{
		{
			name: "Get cloud pocket not found",
			id:   "2",
			sqlFn: func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}

				row := sqlmock.NewRows([]string{
					"id",
					"name",
					"category",
					"amount",
					"goal",
					"currency",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					1,
					"travel",
					"travel",
					1000.0,
					1000.0,
					"THB",
					"2021-09-01T00:00:00Z",
					"2021-09-01T00:00:00Z",
					nil,
				)
				mock.ExpectQuery(cGetOneStmt).WithArgs("1").WillReturnRows(row)
				return db, err
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Get cloud pocket not found because it's deleted",
			id:   "1",
			sqlFn: func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}

				row := sqlmock.NewRows([]string{
					"id",
					"name",
					"category",
					"amount",
					"goal",
					"currency",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					1,
					"travel",
					"travel",
					1000.0,
					1000.0,
					"THB",
					"2021-09-01T00:00:00Z",
					"2021-09-01T00:00:00Z",
					"2021-09-01T00:00:00Z",
				)
				mock.ExpectQuery(cGetOneStmt).WithArgs("1").WillReturnRows(row)
				return db, err
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames(cID)
		c.SetParamValues(tc.id)

		db, err := tc.sqlFn()
		h := New(db)

		assert.NoError(t, err)
		if assert.NoError(t, h.GetOne(c)) {
			assert.Equal(t, tc.wantStatus, rec.Code)
		}
	}
}

func TestGetAll(t *testing.T) {
	e := echo.New()
	var mock sqlmock.Sqlmock
	db, mock, _ := sqlmock.New()
	expect := []GetResponse{
		{
			ID:       1,
			Name:     "travel",
			Category: "Travel",
			Amount:   100,
			Goal:     10_000,
			Currency: THBCurrency,
		},
		{
			ID:       2,
			Name:     "iPhone",
			Category: "Gedged",
			Amount:   100,
			Goal:     40_000,
			Currency: THBCurrency,
		},
	}
	row := sqlmock.NewRows([]string{"ID", "Name", "Category", "Amount", "Goal", "Currency"}).
		AddRow(1, "travel", "Travel", float64(100), float64(10_000), "THB").
		AddRow(2, "iPhone", "Gedged", float64(100), float64(40_000), "THB")

	mock.ExpectPrepare("SELECT id , name , category , amount , goal , currency FROM pockets").ExpectQuery().WithArgs().WillReturnRows(row)

	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := New(db)
	actual := make([]GetResponse, 0)

	assert.NoError(t, h.Get(c))
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expect, actual)

}
