package transaction

import (
	"time"

	cpk "github.com/kkgo-software-engineering/workshop/cloud_pocket"
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

type Transaction struct {
	ID                  int               `json:"id" pg:"pk,unique"`
	Type                TransactionType   `json:"type"`
	Status              TransactionStatus `json:"status"`
	SourcePocketID      int               `json:"sourcePocketId"`
	DestinationPocketID int               `json:"destinationPocketId"`
	Description         string            `json:"description"`
	Amount              float64           `json:"amount"`
	Currency            cpk.Currency      `json:"currency"`
	CreatedAt           time.Time         `json:"createdAt"`
}
