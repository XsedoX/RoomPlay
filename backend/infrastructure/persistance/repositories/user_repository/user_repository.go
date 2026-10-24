package user_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_room_data_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/repositories/common/query"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/repositories/common/specifications"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const getDevicesQuery = `SELECT * FROM devices WHERE user_id = $1`

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

func (repository *UserRepository) Add(ctx context.Context, user *user.User) error {
	userId := user.Id()
	params := []any{
		userId.ToUuid(),
		user.FullName().Name(),
		user.FullName().Surname(),
	}

	// Build VALUES tuples and append deviceFromDb fields to params.
	// Each deviceFromDb contributes 6 columns: id, fingerprint, friendly_name, is_host, type, state
	values := make([]string, 0, len(user.Devices()))
	for i, deviceFromDb := range user.Devices() {
		// parameter indices start after the 5 user params
		base := len(params) + i*5
		// placeholders: ($6,$7,$8,$9,$10,$11), ...
		tuple := fmt.Sprintf("($%d::uuid,$%d,$%d::boolean,$%d::device_type,$%d::device_state)",
			base+1, base+2, base+3, base+4, base+5,
		)
		values = append(values, tuple)

		// append deviceFromDb values in the same order as the tuple
		params = append(params,
			uuid.UUID(deviceFromDb.Id()),
			deviceFromDb.FriendlyName(),
			deviceFromDb.IsHost(),
			deviceFromDb.DeviceType().String(),
			deviceFromDb.State().String(),
		)
	}

	// Compose the CTE + INSERT ... SELECT ... FROM u, (VALUES ...) AS v(...)
	query := `
		WITH "user" AS (
		  INSERT INTO users (id, name, surname)
		  VALUES ($1, $2, $3)
		  RETURNING id
		)
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state)
		SELECT v.id, v.friendly_name, v.is_host, v.type, "user".id, v.state
		FROM "user", (VALUES ` + strings.Join(values, ",") + `) AS v(id, friendly_name, is_host, type, state);`
	_, err := queryer.ExecContext(ctx, query, params...)
	return err
}
