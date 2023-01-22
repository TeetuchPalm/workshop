package pocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Currency string

const (
	THBCurrency Currency = "THB"

	cGetOneStmt = "SELECT * FROM pockets WHERE id = $1"
	cID         = "id"
)

var (
	hErrCloudPocketNotFound = echo.NewHTTPError(http.StatusNotFound,
		"cloud pocket not found")
)

type Err struct {
	Message string `json:"message"`
}

type Pocket struct {
	ID        int        `json:"id" pg:"pk,unique"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	Amount    float64    `json:"amount"`
	Goal      float64    `json:"goal"`
	Currency  Currency   `json:"currency"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

func (h *handler) GetOne(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	pocketID := c.Param(cID)

	pocket := Pocket{}
	var createdAt string
	var updatedAt string
	var deletedAt *string
	if err := h.db.QueryRowContext(ctx, cGetOneStmt, pocketID).Scan(
		&pocket.ID,
		&pocket.Name,
		&pocket.Category,
		&pocket.Amount,
		&pocket.Goal,
		&pocket.Currency,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		logger.Error(fmt.Sprintf("Can not find pocket id: %s", pocketID), zap.Error(err))
		return c.JSON(http.StatusNotFound, hErrCloudPocketNotFound)
	}

	if deletedAt != nil {
		logger.Error(fmt.Sprintf("Can not find pocket id: %s because it's deleted", pocketID))
		return c.JSON(http.StatusNotFound, hErrCloudPocketNotFound)
	}

	pocket.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	pocket.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return c.JSON(http.StatusOK, pocket)
}

func (h *handler) CreatePocket(c echo.Context) error {
	var pocket Pocket
	err := c.Bind(&pocket)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	log.Println(pocket)
	query := "INSERT INTO pockets (name, category, amount, goal, currency, createdat, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	row := h.db.QueryRow(query, pocket.Name, pocket.Category, pocket.Amount, pocket.Goal, pocket.Currency, time.Now(), time.Now())
	err = row.Scan(&pocket.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, pocket)
}

func (h *handler) DeletePocket(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 || err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid id"})
	}
	_, err = h.db.Exec("delete from pockets where id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
