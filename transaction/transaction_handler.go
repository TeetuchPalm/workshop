package transaction

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}
func (h handler) GetTransactionById(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.db.Prepare(`SELECT * FROM transactions WHERE id = $1`)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	row := stmt.QueryRow(id)
	tr := Transaction{}
	err = row.Scan(&tr.ID, &tr.Type, &tr.Status, &tr.SourcePocketID, &tr.DestinationPocketID, &tr.Description, &tr.Amount, &tr.Currency, &tr.CreatedAt)
	switch err {
	case sql.ErrNoRows:
		return c.String(http.StatusNotFound, "get empty row.")
	case nil:
		return c.JSON(http.StatusOK, tr)
	default:
		return c.String(http.StatusInternalServerError, err.Error())
	}
}
