package transaction

import (
	"net/http"
	"time"

	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
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
	Amount              decimal.Decimal   `json:"amount"`
	Currency            pocket.Currency   `json:"currency"`
	CreatedAt           time.Time         `json:"createdAt"`
}

func (h *handler) Transfer(c echo.Context) error {

	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	// get from source pocket
	ctx := c.Request().Context()
	pocket := pocket.Pocket{}
	var amount string
	var goal string
	var createdAt string
	var updatedAt string
	var deletedAt *string
	if err := h.db.QueryRowContext(ctx, "SELECT * FROM pockets WHERE id = $1", t.SourcePocketID).Scan(
		&pocket.ID,
		&pocket.Name,
		&pocket.Category,
		&amount,
		&goal,
		&pocket.Currency,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: err.Error()})
	}

	pocket.Amount, _ = decimal.NewFromString(amount)
	pocket.Goal, _ = decimal.NewFromString(goal)
	pocket.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	pocket.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	query := "SELECT amount FROM pockets WHERE id = $1"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	row := stmt.QueryRow(t.SourcePocketID)
	var am string
	err = row.Scan(&am)

	query = "SELECT amount FROM pockets WHERE id = $1"
	stmt, err = h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	row = stmt.QueryRow(t.DestinationPocketID)
	var amd string
	err = row.Scan(&amd)

	// update to each pocket
	if t.Amount.Cmp(pocket.Amount) == -1 || t.Amount.Cmp(pocket.Amount) == 0 {

		query := "Update pockets set amount = $2, updatedAt = $3 where id = $1"
		stmt, err := h.db.Prepare(query)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}
		sum, err := decimal.NewFromString(am)
		if _, err := stmt.Exec(t.SourcePocketID, sum.Add(t.Amount.Neg()).String(), time.Now()); err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		query = "Update pockets set amount = $2, updatedAt = $3 where id = $1"
		stmt, err = h.db.Prepare(query)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}
		sum2, err := decimal.NewFromString(amd)
		if _, err := stmt.Exec(t.DestinationPocketID, sum2.Add(t.Amount).String(), time.Now()); err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		query = `INSERT INTO transactions (type, status, amount, sourcePocketId, destinationPocketId, description, currency, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
		row = h.db.QueryRow(query, TransferTransactionType, SuccessTransactionStatus, t.Amount, t.SourcePocketID, t.DestinationPocketID, t.Description, t.Currency, time.Now())
		err = row.Scan(&t.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		return c.JSON(http.StatusCreated, t)
	} else {
		query := `INSERT INTO transactions (type, status, amount, sourcePocketId, destinationPocketId, description, currency, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
		row := h.db.QueryRow(query, TransferTransactionType, FailedTransactionStatus, t.Amount.String(), t.SourcePocketID, t.DestinationPocketID, t.Description, t.Currency, time.Now())
		err = row.Scan(&t.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		return c.JSON(http.StatusBadRequest, Err{Message: "not enough money"})
	}

}
