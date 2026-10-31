package internal_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_internal_credentials_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InternalCredentialsRepository struct {
	encrypter i_encrypter.IEncrypter
	uow       *unit_of_work.UnitOfWork
}

func NewInternalCredentialsRepository(encrypter i_encrypter.IEncrypter, uow *unit_of_work.UnitOfWork) *InternalCredentialsRepository {
	return &InternalCredentialsRepository{
		encrypter: encrypter,
		uow:       uow,
	}
}

func (r InternalCredentialsRepository) GetTokenByValue(ctx context.Context, value string) (*internal_credentials.InternalCredentials, error) {
	db := r.uow.GetDB(ctx)
	encryptedRefreshToken := r.encrypter.Hash(value)

	resultDao, err := gorm.G[user_internal_credentials_dao.UserInternalCredentialsDao](db).
		Where("refresh_token = ?", encryptedRefreshToken).
		First(ctx)
	if err != nil {
		return nil, err
	}
	result := resultDao.ToDomain(value)
	return result, nil
}

func (r InternalCredentialsRepository) AssignNewToken(
	ctx context.Context,
	internalCredentials *internal_credentials.InternalCredentials,
) error {
	db := r.uow.GetDB(ctx)
	var dao user_internal_credentials_dao.UserInternalCredentialsDao
	dao.FromDomain(internalCredentials, r.encrypter)
	if _, err := gorm.G[user_internal_credentials_dao.UserInternalCredentialsDao](db).Omit(clause.Associations).
		Updates(ctx, dao); err != nil {
		return err
	}
	return nil
}

func (r InternalCredentialsRepository) RetireTokenByUserSession(ctx context.Context, userSession user_session.UserSession) error {
	db := r.uow.GetDB(ctx)
	if _, err := gorm.G[user_internal_credentials_dao.UserInternalCredentialsDao](db).Where(
		"user_id = ? AND device_id = ?",
		userSession.UserId().ToUuid(),
		userSession.DeviceId().ToUuid(),
	).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (r InternalCredentialsRepository) RetireAllTokensByUserId(ctx context.Context, userId *user_id.UserId) error {
	db := r.uow.GetDB(ctx)
	if _, err := gorm.G[user_internal_credentials_dao.UserInternalCredentialsDao](db).Where(
		"user_id = ?",
		userId.ToUuid(),
	).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}
