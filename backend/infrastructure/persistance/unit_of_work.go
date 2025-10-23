package persistance

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) GetQueryer() IQueryer {
	if uow.tx != nil {
		return uow.db
	}
	return uow.tx
}

func (uow *UnitOfWork) ExecuteTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if uow.tx != nil {
		return fn(ctx)
	}
	tx, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	uow.tx = tx
	defer func() {
		// clear tx pointer regardless of outcome to avoid reuse
		defer func() { uow.tx = nil }()

		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		// commit error (if any) should be returned to caller
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			return
		}
	}()

	err = fn(ctx)
	return err
}

func (uow *UnitOfWork) ExecuteRead(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
