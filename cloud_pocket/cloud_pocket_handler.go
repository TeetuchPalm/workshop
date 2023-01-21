package cloud_pocket

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
func (h handler) Get(c echo.Context) error {
	stmt, err := h.db.Prepare("SELECT * FROM pocket")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	pocket := []CloudPocket{}

	for rows.Next() {
		ex := CloudPocket{}
		err := rows.Scan(&ex.ID, &ex.Name, &ex.Category, &ex.Amount, &ex.Goal, &ex.Currency)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		pocket = append(pocket, ex)
	}

	return c.JSON(http.StatusOK, pocket)
}
