package user_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_room_data_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/repositories/common/query"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/repositories/common/specifications"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	unitOfWork *unit_of_work.UnitOfWork
}

func NewUserRepository(uow *unit_of_work.UnitOfWork) *UserRepository {
	return &UserRepository{
		unitOfWork: uow,
	}
}

func (repository *UserRepository) GetUserById(ctx context.Context, id user_id.UserId) (*user.User, error) {
	db := repository.unitOfWork.GetDB(ctx)

	value, err := gorm.G[user_dao.UserDao](db).
		Where("id = ?", id.ToUuid()).
		Preload("Devices", nil).
		Preload("RoomData", nil).
		First(ctx)
	if err != nil {
		return nil, err
	}
	result := value.ToDomain()

	return result, nil
}

func (repository *UserRepository) GetUserByExternalId(ctx context.Context, externalId string) (*user.User, error) {
	db := repository.unitOfWork.GetDB(ctx)

	value, err := gorm.G[user_dao.UserDao](db).
		Where("external_id = ?", externalId).
		Preload("Devices", nil).
		Preload("RoomData", nil).
		First(ctx)
	if err != nil {
		return nil, err
	}
	result := value.ToDomain()

	return result, nil
}

func (repository *UserRepository) CheckIfUserExistByExternalId(ctx context.Context, externalId string) (bool, error) {
	db := repository.unitOfWork.GetDB(ctx)
	response, err := query.Exists[user_dao.UserDao](
		ctx,
		db,
		specifications.WithExternalCredentials,
		specifications.WithExternalID(externalId),
	)
	if err != nil {
		return false, err
	}
	return response, nil
}

func (repository *UserRepository) Update(ctx context.Context, user *user.User) error {
	db := repository.unitOfWork.GetDB(ctx)
	var roomData *user_room_data_dao.UserRoomDataDao
	if user.RoomId() != nil && user.Role() != nil && user.BoostUsedAtUtc() != nil {
		roomData = &user_room_data_dao.UserRoomDataDao{
			BoostUsedAtUtc: user.BoostUsedAtUtc(),
		}
	}
	userModel := &user_dao.UserDao{
		Id:      user.Id().ToUuid(),
		Name:    user.FullName().Name(),
		Surname: user.FullName().Surname(),
	}
	if _, err := gorm.G[user_dao.UserDao](db).Where(
		"id = ?",
		user.Id().ToUuid(),
	).Omit(clause.Associations).
		Updates(ctx, *userModel); err != nil {
		return err
	}
	if _, err := gorm.G[device_dao.DeviceDao](
		db.Scopes(specifications.WithUserID(user.Id())),
	).Delete(ctx); err != nil {
		return err
	}

	if len(user.Devices()) > 0 {
		devices := make([]device_dao.DeviceDao, len(user.Devices()))
		for i := range user.Devices() {
			var deviceDao device_dao.DeviceDao
			userDevice := user.Devices()[i]
			deviceDao.FromDomain(&userDevice, user.Id())
			devices[i] = deviceDao
		}

		if err := gorm.G[device_dao.DeviceDao](db).CreateInBatches(
			ctx,
			&devices,
			100,
		); err != nil {
			return err
		}
	}
	if roomData != nil {
		roomData.UserId = user.Id().ToUuid()
		err := gorm.G[user_room_data_dao.UserRoomDataDao](db, clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "room_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"role", "boost_used_at_utc"}),
		}).Create(ctx, roomData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository *UserRepository) Add(ctx context.Context, user *user.User) (*user_id.UserId, error) {
	db := repository.unitOfWork.GetDB(ctx)
	userId := user.Id()
	userModel := &user_dao.UserDao{
		Id:      userId.ToUuid(),
		Name:    user.FullName().Name(),
		Surname: user.FullName().Surname(),
	}
	if err := gorm.G[user_dao.UserDao](db).Omit(clause.Associations).Create(ctx, userModel); err != nil {
		return nil, err
	}

	if _, err := gorm.G[device_dao.DeviceDao](
		db.Scopes(specifications.WithUserID(user.Id())),
	).Delete(ctx); err != nil {
		return nil, err
	}
	if len(user.Devices()) > 0 {
		devices := make([]device_dao.DeviceDao, len(user.Devices()))
		for i := range user.Devices() {
			var deviceDao device_dao.DeviceDao
			userDevice := user.Devices()[i]
			deviceDao.FromDomain(&userDevice, user.Id())
			devices[i] = deviceDao
		}

		if err := gorm.G[device_dao.DeviceDao](db).CreateInBatches(
			ctx,
			&devices,
			100,
		); err != nil {
			return nil, err
		}
	}
	return &userId, nil
}
