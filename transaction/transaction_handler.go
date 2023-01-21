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
func (h handler) GetTransactionByPocketId(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.db.Prepare(`SELECT * FROM transactions WHERE sourcePocketId = $1 OR destinationPocketId = $1 ORDER BY id`)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	trs := []Transaction{}
	tr := Transaction{}
	for rows.Next() {
		err = rows.Scan(&tr.ID, &tr.Type, &tr.Status, &tr.SourcePocketID, &tr.DestinationPocketID, &tr.Description, &tr.Amount, &tr.Currency, &tr.CreatedAt)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		trs = append(trs, tr)

	}
	return c.JSON(http.StatusOK, trs)
}
