package cloud_pocket

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Currency string

const (
	THBCurrency Currency = "THB"

	cGetStmt = "SELECT * FROM cloud_pockets WHERE id = $1"
	cID      = "id"
)

var (
	hErrCloudPocketNotFound = echo.NewHTTPError(http.StatusNotFound,
		"cloud pocket not found")
)

type CloudPocket struct {
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

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) Get(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	pocketID := c.Param(cID)

	pocket := CloudPocket{}
	var createdAt string
	var updatedAt string
	var deletedAt *string
	if err := h.db.QueryRowContext(ctx, cGetStmt, pocketID).Scan(
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
