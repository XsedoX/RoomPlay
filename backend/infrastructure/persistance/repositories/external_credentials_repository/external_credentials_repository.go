package external_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_external_credentials_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExternalCredentialsRepository struct {
	encrypter i_encrypter.IEncrypter
	uow       *unit_of_work.UnitOfWork
}

func NewExternalCredentialsRepository(encrypter i_encrypter.IEncrypter, uow *unit_of_work.UnitOfWork) *ExternalCredentialsRepository {
	return &ExternalCredentialsRepository{
		encrypter: encrypter,
		uow:       uow,
	}
}

func (repo *ExternalCredentialsRepository) UpsertExternalCredentials(ctx context.Context, credentials *external_credentials.ExternalCredentials) error {
	db := repo.uow.GetDB(ctx)
	var dao user_external_credentials_dao.UserExternalCredentialsDao
	if err := dao.FromDomain(credentials, repo.encrypter); err != nil {
		return err
	}
	if _, err := gorm.G[user_external_credentials_dao.UserExternalCredentialsDao](db).Omit(clause.Associations).Updates(ctx, dao); err != nil {
		return nil
	}
	return nil
}

func (repo *ExternalCredentialsRepository) GetExternalCredentialsByUserId(ctx context.Context, userId user_id.UserId) (*external_credentials.ExternalCredentials, error) {
	db := repo.uow.GetDB(ctx)
	dao, err := gorm.G[user_external_credentials_dao.UserExternalCredentialsDao](db).
		Where("user_id = ?", userId.ToUuid()).First(ctx)
	if err != nil {
		return nil, err
	}
	result, err := dao.ToDomain(repo.encrypter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
