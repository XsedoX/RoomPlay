package persistance

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/user"
)

type UserRepository struct {
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

func (repository *UserRepository) UpdateUser(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {
	_, err := queryer.ExecContext(ctx, `
		UPDATE users 
		SET name=$1, surname=$2, room_id=$3::uuid
		WHERE id=$4::uuid`,
		user.Name(), user.Surname(), uuid.UUID(*user.RoomId()), uuid.UUID(user.Id()))
	if err != nil {
		return err
	}

	values := make([]string, 0, len(user.Devices()))
	var params []interface{}
	base := 0
	for i, device := range user.Devices() {
		base = base + i*4
		tuple := fmt.Sprintf("($%d::uuid,$%d,$%d,$%d::device_state)",
			base+1, base+2, base+3, base+4,
		)
		values = append(values, tuple)
		params = append(params,
			uuid.UUID(device.Id()),
			device.FriendlyName(),
			device.IsHost(),
			device.State().String())
	}
	query := `
		UPDATE devices AS d
		SET d.friendly_name=c.friendly_name, d.is_host=c.is_host, d.state=c.state FROM
		(VALUES + strings.Join(values, ",") + ) AS c(id, friendly_name, is_host, state)
		WHERE d.id = c.id;`

	_, err = queryer.ExecContext(ctx, query, params...)

	return err
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repository *UserRepository) AddUser(ctx context.Context, user *user.User, queryer contracts.IQueryer) error {

	params := []interface{}{
		uuid.UUID(user.Id()),
		user.ExternalId(),
		user.Name(),
		user.Surname(),
		uuid.UUID(*user.RoomId()),
	}

	// If no devices, do a simple single INSERT
	if len(user.Devices()) == 0 {
		_, err := queryer.ExecContext(ctx,
			"INSERT INTO users (id, external_id, name, surname, room_id) VALUES ($1::uuid, $2, $3, $4, $5::uuid);",
			params...,
		)
		return err
	}

	// Build VALUES tuples and append device fields to params.
	// Each device contributes 6 columns: id, fingerprint, friendly_name, is_host, type, state
	values := make([]string, 0, len(user.Devices()))
	for i, device := range user.Devices() {
		// parameter indices start after the 5 user params
		base := len(params) + i*5
		// placeholders: ($6,$7,$8,$9,$10,$11), ...
		tuple := fmt.Sprintf("($%d::uuid,$%d,$%d::boolean,$%d::device_type,$%d::device_state)",
			base+1, base+2, base+3, base+4, base+5,
		)
		values = append(values, tuple)

		// append device values in the same order as the tuple
		params = append(params,
			uuid.UUID(device.Id()),
			device.FriendlyName(),
			device.IsHost(),
			device.DeviceType().String(),
			device.State().String(),
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
