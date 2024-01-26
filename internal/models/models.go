// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package models

import (
	"database/sql"

	types "finance/internal/types"
)

type Transaction struct {
	ID          int64
	Date        types.Date
	Code        sql.NullString
	Description string
	Amount      int64
	Balance     int64
}
