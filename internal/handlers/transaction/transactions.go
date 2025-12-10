package transaction

import (
	"context"
	"database/sql"
)

type txKey struct{}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	GetTx(ctx context.Context) *sql.Tx
}

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	err = fn(txCtx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tm *transactionManager) GetTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txKey{}).(*sql.Tx)
	return tx
}
