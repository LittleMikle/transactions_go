package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
)

type RequestInterface interface {
	Validate(db *pgx.Conn) error
}

type Deposit struct {
	Receiver uint32  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func (d Deposit) Validate(db *pgx.Conn) error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Receiver, validation.Required, validation.By(validateExists(db, "wallets"))),
		validation.Field(
			&d.Amount,
			validation.Required,
			validation.Min(0.01),
			validation.By(validateDeposit(db, "wallets", "amount", d.Receiver)),
		),
	)
}

type Transfer struct {
	Sender   uint32  `json:"sender"`
	Receiver uint32  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func (t Transfer) Validate(db *pgx.Conn) error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Sender, validation.Required, validation.By(validateExists(db, "wallets"))),
		validation.Field(&t.Receiver, validation.Required, validation.By(validateExists(db, "wallets"))),
		validation.Field(
			&t.Amount,
			validation.Required,
			validation.Min(0.01),
			validation.By(validateDeposit(db, "wallets", "amount", t.Receiver)),
			validation.By(validateWithdrawal(db, "wallets", "amount", t.Sender)),
		),
	)
}
