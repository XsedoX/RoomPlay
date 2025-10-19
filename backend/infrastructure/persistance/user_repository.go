package persistance

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/domain/entities"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) Create(ctx context.Context, user *entities.User) error {
	queryer := GetQueryerFromContext(ctx, repository.db)

	_, err := queryer.QueryxContext(ctx,
		"INSERT INTO users (id, external_id, name, surname, room_id) VALUES ($1, $2, $3, $4, $5);",
		uuid.UUID(user.Id()),
		user.ExternalId(),
		user.Name(),
		user.Surname(),
		user.RoomId())
	if err != nil {
		return err
	}
	return nil
}
