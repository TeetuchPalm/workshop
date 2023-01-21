//go:build unit

package cloud_pocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const endpoint = "/cloud_pocket"

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
				mock.ExpectQuery(cGetStmt).WithArgs("1").WillReturnRows(row)
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
		req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(tc.reqBody))
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
				mock.ExpectQuery(cGetStmt).WithArgs("1").WillReturnRows(row)
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
				mock.ExpectQuery(cGetStmt).WithArgs("1").WillReturnRows(row)
				return db, err
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(tc.reqBody))
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
