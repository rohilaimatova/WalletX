package repository

import (
	"WalletX/pkg/errs"
	"database/sql"
)

// Перевод ошибок БД в наши ошибки
func translateError(err error) error {
	if err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return errs.ErrUserNotFound
	} else {
		return errs.ErrInternal
	}
}
