package persistance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

const getDevicesQuery = `SELECT * FROM devices WHERE user_id = $1`

type userDbDto struct {
	Id         uuid.UUID  `db:"id"`
	ExternalId string     `db:"external_id"`
	Name       string     `db:"name"`
	Surname    string     `db:"surname"`
	RoomId     *uuid.UUID `db:"room_id"`
	Role       *string    `db:"role"`
}
type deviceDbDto struct {
	Id                uuid.UUID `db:"id"`
	FriendlyName      string    `db:"friendly_name"`
	IsHost            bool      `db:"is_host"`
	Type              string    `db:"type"`
	UserId            uuid.UUID `db:"user_id"`
	State             string    `db:"state"`
	LastLoggedInAtUtc time.Time `db:"last_logged_in_at_utc"`
}
type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
func (repository *UserRepository) GetUserById(ctx context.Context, id shared.UserId, queryer contracts.IQueryer) (*user.User, error) {
	var userDb userDbDto
	err := queryer.GetContext(ctx,
		&userDb,
		`SELECT u.*, ur.role FROM users u 
         	   LEFT JOIN users_roles ur ON ur.user_id = u.id 
         	   WHERE id = $1`,
		uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	var devicesDb []*deviceDbDto
	err = queryer.SelectContext(ctx,
		&devicesDb,
		getDevicesQuery,
		userDb.Id)
	if err != nil {
		return nil, err
	}

	return parseUser(&userDb, devicesDb), nil
}
func (repository *UserRepository) GetUserByExternalId(ctx context.Context, externalId string, queryer contracts.IQueryer) (*user.User, error) {
	var userDb userDbDto
	err := queryer.GetContext(ctx,
		&userDb,
		`SELECT u.*, ur.role FROM users u 
         	   LEFT JOIN users_roles ur ON ur.user_id = u.id 
         	   WHERE external_id = $1`,
		externalId)
	if err != nil {
		return nil, err
	}

	var devicesDb []*deviceDbDto
	err = queryer.SelectContext(ctx,
		&devicesDb,
		getDevicesQuery,
		userDb.Id)
	if err != nil {
		return nil, err
	}
	return parseUser(&userDb, devicesDb), nil
}
func (repository *UserRepository) CheckIfUserExistByExternalId(ctx context.Context, externalId string, queryer contracts.IQueryer) bool {
	var response bool
	err := queryer.GetContext(ctx, &response, `
		SELECT CASE 
		    WHEN EXISTS (
		        SELECT 1
		        FROM users
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
func (repository *UserRepository) Update(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {
	var roomId *uuid.UUID
	if user.RoomId() != nil {
		id := uuid.UUID(*user.RoomId())
		roomId = &id
	}
	_, err := queryer.ExecContext(ctx, `
		UPDATE users 
		SET name=$1, surname=$2, room_id=$3::uuid
		WHERE id=$4::uuid`,
		user.FullName().Name(), user.FullName().Surname(), roomId, uuid.UUID(user.Id()))
	if err != nil {
		return err
	}

	if user.Role() != nil && user.RoomId() != nil {
		_, err = queryer.ExecContext(ctx,
			`UPDATE users_roles SET role=$1 WHERE user_id = $2::uuid AND room_id = $3::uuid`,
			user.Role().String(),
			uuid.UUID(user.Id()),
			uuid.UUID(*user.RoomId()))
		if err != nil {
			return err
		}
	}

	values := make([]string, 0, len(user.Devices()))
	var params []interface{}
	for i, deviceToUpdate := range user.Devices() {
		base := i * 7
		tuple := fmt.Sprintf("($%d::uuid, $%d, $%d, $%d::device_type, $%d::uuid, $%d::device_state, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7,
		)
		values = append(values, tuple)
		params = append(params,
			uuid.UUID(deviceToUpdate.Id()),
			deviceToUpdate.FriendlyName(),
			deviceToUpdate.IsHost(),
			deviceToUpdate.DeviceType().String(),
			uuid.UUID(user.Id()),
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

	values = make([]string, 0, len(user.Devices()))
	params1 := make([]interface{}, 0, len(user.Devices())+1)
	params1 = append(params1, uuid.UUID(user.Id()))
	for i, deviceToUpdate := range user.Devices() {
		values = append(values, fmt.Sprintf("$%d", i+2))
		params1 = append(params1, uuid.UUID(deviceToUpdate.Id()))
	}

	deleteQuery := `DELETE FROM devices WHERE user_id=$1`
	if len(values) > 0 {
		deleteQuery += ` AND id NOT IN (` + strings.Join(values, ",") + `)`
	}

	_, err = queryer.ExecContext(ctx, deleteQuery, params1...)

	return err
}
func (repository *UserRepository) Add(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {
	var roomId *uuid.UUID
	if user.RoomId() != nil {
		id := uuid.UUID(*user.RoomId())
		roomId = &id
	}
	params := []interface{}{
		uuid.UUID(user.Id()),
		user.ExternalId(),
		user.FullName().Name(),
		user.FullName().Surname(),
		roomId,
	}

	// If no devices, do a simple single INSERT
	if len(user.Devices()) == 0 {
		_, err := queryer.ExecContext(ctx,
			"INSERT INTO users (id, external_id, name, surname, room_id) VALUES ($1::uuid, $2, $3, $4, $5::uuid);",
			params...,
		)
		return err
	}

	if user.RoomId() != nil && user.Role() != nil {
		_, err := queryer.ExecContext(ctx,
			"INSERT INTO users_roles (room_id, user_id, role) VALUES ($1, $2, $3::user_role)",
			uuid.UUID(*user.RoomId()),
			uuid.UUID(user.Id()),
			user.Role().String(),
		)
		return err
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
		  INSERT INTO users (id, external_id, name, surname, room_id)
		  VALUES ($1, $2, $3, $4, $5)
		  RETURNING id
		)
		INSERT INTO devices (id, friendly_name, is_host, type, user_id, state)
		SELECT v.id, v.friendly_name, v.is_host, v.type, "user".id, v.state
		FROM "user", (VALUES ` + strings.Join(values, ",") + `) AS v(id, friendly_name, is_host, type, state);`
	_, err := queryer.ExecContext(ctx, query, params...)
	return err
}
func parseUser(userDb *userDbDto, devicesDb []*deviceDbDto) *user.User {
	var devices []device.Device
	for _, deviceDb := range devicesDb {
		deviceResult := device.HydrateDevice(
			shared.DeviceId(deviceDb.Id),
			deviceDb.FriendlyName,
			*device.ParseType(deviceDb.Type),
			deviceDb.IsHost,
			device.ParseState(deviceDb.State),
			deviceDb.LastLoggedInAtUtc,
		)
		devices = append(devices,
			*deviceResult)
	}

	return user.HydrateUser(shared.UserId(userDb.Id),
		userDb.ExternalId,
		userDb.Name,
		userDb.Surname,
		user.ParseRole(userDb.Role),
		(*shared.RoomId)(userDb.RoomId),
		devices)
}
