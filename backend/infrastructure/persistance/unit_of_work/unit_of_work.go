package unit_of_work

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db *sqlx.DB
}

type ctxKey struct{}

func withTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func txFromCtx(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(ctxKey{}).(*sqlx.Tx)
	return tx
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) GetQueryer(ctx context.Context) i_queryer.IQueryer {
	if tx := txFromCtx(ctx); tx != nil {
		return tx
	}
	return uow.db
}

func (uow *UnitOfWork) ExecuteTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if txFromCtx(ctx) != nil {
		return fn(ctx)
	}
	tx, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	var commitErr error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if commitErr != nil {
			_ = tx.Rollback()
			return
		}
		if cErr := tx.Commit(); cErr != nil {
			commitErr = cErr
			_ = tx.Rollback()
		}
	}()

	txCtx := withTx(ctx, tx)
	commitErr = fn(txCtx)
	return commitErr
}

func (uow *UnitOfWork) ExecuteRead(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
