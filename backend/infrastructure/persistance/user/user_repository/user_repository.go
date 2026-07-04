package user_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/domain/user/device"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_state"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_type"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_dao"
	"github.com/google/uuid"
)

const getDevicesQuery = `SELECT * FROM devices WHERE user_id = $1`

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repository *UserRepository) GetUserById(ctx context.Context, id user_id.UserId, queryer i_queryer.IQueryer) (*user.User, error) {
	var userDb user_dao.UserDao
	err := queryer.GetContext(ctx,
		&userDb,
		`SELECT u.*, ur.role, ur.boost_used_at_utc, ur.room_id FROM users u 
         	   LEFT JOIN users_room_data ur ON ur.user_id = u.id 
         	   WHERE id = $1`,
		id.ToUuid())
	if err != nil {
		return nil, err
	}

	devicesDb := make([]device_dao.DeviceDao, 0)
	err = queryer.SelectContext(ctx,
		&devicesDb,
		getDevicesQuery,
		userDb.Id)
	if err != nil {
		return nil, err
	}

	return parseUser(&userDb, &devicesDb), nil
}

func (repository *UserRepository) GetUserByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) (*user.User, error) {
	var userDb user_dao.UserDao
	err := queryer.GetContext(ctx,
		&userDb,
		`SELECT u.*, ur.role, ur.boost_used_at_utc, ur.room_id FROM users u 
         	   LEFT JOIN users_room_data ur ON ur.user_id = u.id 
		LEFT JOIN users_external_credentials uec ON uec.user_id = u.id
         	   WHERE uec.external_id = $1`,
		externalId)
	if err != nil {
		return nil, err
	}

	devicesDb := make([]device_dao.DeviceDao, 0)
	err = queryer.SelectContext(ctx,
		&devicesDb,
		getDevicesQuery,
		userDb.Id)
	if err != nil {
		return nil, err
	}
	return parseUser(&userDb, &devicesDb), nil
}

func (repository *UserRepository) CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer i_queryer.IQueryer) bool {
	var response bool
	err := queryer.GetContext(ctx, &response, `
		SELECT CASE 
		    WHEN EXISTS (
		        SELECT 1
		        FROM users_external_credentials
		        WHERE external_id=$1
			)
		    THEN true 
		    ELSE false
		END`, externalId)
	if err != nil {
		return false
	}
	return response
}

func (repository *UserRepository) Update(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error {
	userId := user.Id()
	_, err := queryer.ExecContext(ctx, `
		UPDATE users 
		SET name=$1, surname=$2
		WHERE id=$3::uuid`,
		user.FullName().Name(),
		user.FullName().Surname(),
		userId.ToUuid())
	if err != nil {
		return err
	}

	if user.Role() != nil && user.RoomId() != nil && user.BoostUsedAtUtc() != nil {
		_, err = queryer.ExecContext(ctx,
			`INSERT INTO users_room_data (room_id, user_id, role, boost_used_at_utc) 
				   VALUES ($1, $2, $3, $4)
				   ON CONFLICT (room_id, user_id) DO UPDATE SET role=$3::user_role, boost_used_at_utc=$4;`,
			user.RoomId().ToUuid(),
			userId.ToUuid(),
			user.Role().String(),
			user.BoostUsedAtUtc())
		if err != nil {
			return err
		}
	}

	if len(user.Devices()) > 0 {
		values := make([]string, 0, len(user.Devices()))
		var params []any
		for i, deviceToUpdate := range user.Devices() {
			base := i * 7
			tuple := fmt.Sprintf("($%d::uuid, $%d, $%d, $%d::device_type, $%d::uuid, $%d::device_state, $%d)",
				base+1, base+2, base+3, base+4, base+5, base+6, base+7,
			)
			values = append(values, tuple)
			deviceId := deviceToUpdate.Id()
			params = append(params,
				deviceId.ToUuid(),
				deviceToUpdate.FriendlyName(),
				deviceToUpdate.IsHost(),
				deviceToUpdate.DeviceType().String(),
				userId.ToUuid(),
				deviceToUpdate.State().String(),
				deviceToUpdate.LastLoggedInUtc(),
			)
		}
		query := `
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state, last_logged_in_at_utc)
		VALUES ` + strings.Join(values, ",") + ` 
		ON CONFLICT (id, user_id) DO UPDATE
		SET
			friendly_name = EXCLUDED.friendly_name,
			is_host = EXCLUDED.is_host,
			type = EXCLUDED.type,
			state = EXCLUDED.state,
			last_logged_in_at_utc = EXCLUDED.last_logged_in_at_utc;`

		_, err = queryer.ExecContext(ctx, query, params...)
		if err != nil {
			return err
		}
	}

	if len(user.Devices()) == 0 {
		_, err = queryer.ExecContext(ctx, `DELETE FROM devices WHERE user_id=$1`, userId.ToUuid())
		return err
	}

	values := make([]string, 0, len(user.Devices()))
	params1 := make([]any, 0, len(user.Devices())+1)
	params1 = append(params1, userId.ToUuid())
	for i, deviceToUpdate := range user.Devices() {
		values = append(values, fmt.Sprintf("$%d", i+2))
		deviceId := deviceToUpdate.Id()
		params1 = append(params1, deviceId.ToUuid())
	}

	deleteQuery := `DELETE FROM devices WHERE user_id=$1
					AND id NOT IN (` + strings.Join(values, ",") + `);`

	_, err = queryer.ExecContext(ctx, deleteQuery, params1...)

	return err
}

func (repository *UserRepository) Add(ctx context.Context, user *user.User, queryer i_queryer.IQueryer) error {
	userId := user.Id()
	params := []interface{}{
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

func parseUser(userDb *user_dao.UserDao, devicesDb *[]device_dao.DeviceDao) *user.User {
	var devices []device.Device
	for _, deviceDb := range *devicesDb {
		deviceResult := device.HydrateDevice(
			device_id.DeviceId(deviceDb.Id),
			deviceDb.FriendlyName,
			*device_type.ParseDeviceType(&deviceDb.Type),
			deviceDb.IsHost,
			*device_state.ParseDeviceState(&deviceDb.State),
			deviceDb.LastLoggedInAtUtc,
		)
		devices = append(devices,
			*deviceResult)
	}

	var roomId *room_id.RoomId
	if userDb.RoomId == nil {
		roomId = (*room_id.RoomId)(nil)
	} else {
		roomIdConcrete := room_id.RoomId(*userDb.RoomId)
		roomId = &roomIdConcrete
	}
	return user.HydrateUser(user_id.UserId(userDb.Id),
		userDb.Name,
		userDb.Surname,
		user_role.ParseUserRole(userDb.Role),
		roomId,
		devices,
		userDb.BoostUsedAtUtc)
}
