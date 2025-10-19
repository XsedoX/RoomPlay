package persistance

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const UnitOfWorkCtxKey = "UnitOfWork"

type UnitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	transaction, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()
			panic(p)
		}
	}()

	transactionCtx := context.WithValue(ctx, UnitOfWorkCtxKey, transaction)

	if err := fn(transactionCtx); err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w", rbErr)
		}
		return err
	}
	return transaction.Commit()
}

func GetQueryerFromContext(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(UnitOfWorkCtxKey).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
