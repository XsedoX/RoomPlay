package persistance

import (
	"context"

	"github.com/google/uuid"
	"xsedox.com/main/domain/entities"
)

type UserRepository struct {
	queryer IQueryer
}

func NewUserRepository(q IQueryer) *UserRepository {
	return &UserRepository{queryer: q}
}

func (repository *UserRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := repository.queryer.ExecContext(ctx,
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
