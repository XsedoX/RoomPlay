package validation

import (
	"github.com/go-playground/validator/v10"
	"xsedox.com/main/domain/entities"
)

var Validator *validator.Validate

func Initialize() error {
	Validator = validator.New(validator.WithRequiredStructEnabled())
	err := Validator.RegisterValidation("user_role_validation", validateUserRoleEnum)
	if err != nil {
		return err
	}
	err = Validator.RegisterValidation("device_type_validation", validateDeviceType)
	if err != nil {
		return err
	}
	return nil
}

func validateUserRoleEnum(fl validator.FieldLevel) bool {
	role := entities.UserRole(fl.Field().String())
	if role == entities.HOST || role == entities.USER {
		return true
	}
	return false
}
func validateDeviceType(fl validator.FieldLevel) bool {
	deviceType := entities.DeviceType(fl.Field().String())
	if deviceType == entities.MOBILE || deviceType == entities.COMPUTER {
		return true
	}
	return false
}
