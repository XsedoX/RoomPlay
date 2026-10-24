// Package unit_of_work provides a Unit of Work implementation for managing database transactions using GORM.
package unit_of_work

import (
	"context"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
}

type transactionKey struct{}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) ExecuteTransaction(ctx context.Context, fn func(context.Context) error) error {
	if txFromCtx(ctx) != nil {
		return fn(ctx)
	}

	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey{}, tx)
		return fn(txCtx)
	})
}

func (uow *UnitOfWork) GetDB(ctx context.Context) *gorm.DB {
	if tx := txFromCtx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return uow.db.WithContext(ctx)
}

func txFromCtx(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(transactionKey{}).(*gorm.DB)
	return tx
}
