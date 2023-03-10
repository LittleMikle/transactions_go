package request

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
)

const (
	numericPrecision = 13
	maxNumeric       = 9999999999999.99
)

func validateExists(db *pgx.Conn, table string) validation.RuleFunc {
	return func(value interface{}) error {
		var (
			col int
			id  = value.(uint32)
			qry = fmt.Sprintf("SELECT id FROM %s WHERE id = $1", table)
		)
		if err := db.QueryRow(context.Background(), qry, id).Scan(&col); err != nil {
			return errors.New("wallet not found")
		}

		return nil
	}
}

func validateDeposit(db *pgx.Conn, table string, col string, id uint32) validation.RuleFunc {
	return func(value interface{}) error {
		var (
			amnt float64
			qry  = fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", col, table)
		)
		row := db.QueryRow(context.Background(), qry, id)
		if err := row.Scan(&amnt); err != nil {
			return nil
		}

		if amnt+value.(float64) > maxNumeric {
			return fmt.Errorf("too much money! Result amount must be no greater than 1e+%v", numericPrecision)
		}

		return nil
	}
}

func validateWithdrawal(db *pgx.Conn, table string, col string, id uint32) validation.RuleFunc {
	return func(value interface{}) error {
		var (
			amnt float64
			qry  = fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", col, table)
		)
		row := db.QueryRow(context.Background(), qry, id)
		if err := row.Scan(&amnt); err != nil {
			return nil
		}

		if value.(float64) > amnt {
			return errors.New("insufficient funds")
		}

		return nil
	}
}
