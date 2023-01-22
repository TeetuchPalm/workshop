package transaction

import (
	"net/http"
	"time"

	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
)

type (
	TransactionType   string
	TransactionStatus string
)

const (
	DepositTransactionType  TransactionType = "deposit"
	WithdrawTransactionType TransactionType = "withdraw"
	TransferTransactionType TransactionType = "transfer"

	SuccessTransactionStatus TransactionStatus = "success"
	FailedTransactionStatus  TransactionStatus = "failed"
)

type Err struct {
	Message string `json:"message"`
}

type Transaction struct {
	ID                  int               `json:"id" pg:"pk,unique"`
	Type                TransactionType   `json:"type"`
	Status              TransactionStatus `json:"status"`
	SourcePocketID      int               `json:"sourcePocketId"`
	DestinationPocketID int               `json:"destinationPocketId"`
	Description         string            `json:"description"`
	Amount              float64           `json:"amount"`
	Currency            pocket.Currency   `json:"currency"`
	CreatedAt           time.Time         `json:"createdAt"`
}

func (h *handler) Transfer(c echo.Context) error {

	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	query := "Update pockets set amount = amount-$2, updatedAt = $3 where id = $1"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if _, err := stmt.Exec(t.SourcePocketID, t.Amount, time.Now()); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	query = "Update pockets set amount = amount+$2, updatedAt = $3 where id = $1"
	stmt, err = h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if _, err := stmt.Exec(t.DestinationPocketID, t.Amount, time.Now()); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	query = `INSERT INTO transactions (type, status, amount, sourcePocketId, destinationPocketId, description, currency, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	row := h.db.QueryRow(query, TransferTransactionType, SuccessTransactionStatus, t.Amount, t.SourcePocketID, t.DestinationPocketID, t.Description, t.Currency, time.Now())
	err = row.Scan(&t.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, t)
}