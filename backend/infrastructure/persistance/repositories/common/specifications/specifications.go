package specifications

import (
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/repositories/common/query"
	"gorm.io/gorm"
)

func WithExternalCredentials(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN users_external_credentials ON users_external_credentials.user_id = users.id")
}

func WithExternalID(externalID string) query.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("users_external_credentials.external_id = ?", externalID)
	}
}

func WithUserID(userID user_id.UserId) query.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("users.id = ?", userID.ToUuid())
	}
}
