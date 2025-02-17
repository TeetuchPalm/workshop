package pocket

import (
	"database/sql"
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
	hErrCloudPocketNotFound = echo.NewHTTPError(
		http.StatusNotFound,
		"cloud pocket not found",
	)
	hErrCloudPocketDeleteFailed = echo.NewHTTPError(
		http.StatusNotFound,
		"cloud pocket delete failed",
	)	
	hErrGetAllFailed = echo.NewHTTPError(
		http.StatusInternalServerError,
		"กรุณาติดต่อผู้ดูแลระบบ (RefCode:12)",
	)
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

type GetResponse struct {
	ID       int      `json:"id" pg:"pk,unique"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Amount   float64  `json:"amount"`
	Goal     float64  `json:"goal"`
	Currency Currency `json:"currency"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
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
	logger := mlog.L(c)
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 || err != nil {
		logger.Error(fmt.Sprintf("delete pocket invalid id %d", id))
		return c.JSON(http.StatusBadRequest, hErrCloudPocketDeleteFailed)
	}

	var isDeleted int
	err = h.db.QueryRowContext(c.Request().Context(),
		"select 1 from pockets where id = $1 and deletedat is null", id).Scan(&isDeleted)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete pocket. %s", err.Error()))
		return c.JSON(http.StatusBadRequest, hErrCloudPocketDeleteFailed)
	}

	_, err = h.db.Exec("update pockets set deletedat = $2 where id = $1", id, time.Now())
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete pocket. %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, hErrCloudPocketDeleteFailed)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) Get(c echo.Context) error {
	logger := mlog.L(c)

	stmt, err := h.db.Prepare("SELECT id , name , category , amount , goal , currency FROM pockets")
	if err != nil {
		logger.Error("Can not prepare sql", zap.Error(err))
		return hErrGetAllFailed
	}

	rows, err := stmt.Query()
	if err != nil {
		logger.Error("Can not query sql", zap.Error(err))
		return hErrGetAllFailed
	}

	pocket := []GetResponse{}

	for rows.Next() {
		res := GetResponse{}
		err := rows.Scan(&res.ID, &res.Name, &res.Category, &res.Amount, &res.Goal, &res.Currency)
		if err != nil {
			logger.Error("Can not scan rows", zap.Error(err))
			return hErrGetAllFailed
		}
		pocket = append(pocket, res)
	}

	return c.JSON(http.StatusOK, pocket)
}
