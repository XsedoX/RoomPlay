package validation

import (
	"github.com/go-playground/validator/v10"
	"xsedox.com/main/domain/entities"
)

var Validate *validator.Validate

func Initialize() error {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	err := Validate.RegisterValidation("user_role_validation", ValidateUserRoleEnum)
	if err != nil {
		return err
	}
	err = Validate.RegisterValidation("device_type_validation", ValidateDeviceType)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserRoleEnum(fl validator.FieldLevel) bool {
	role := entities.UserRole(fl.Field().String())
	if role == entities.HOST || role == entities.USER {
		return true
	}
	return false
}
func ValidateDeviceType(fl validator.FieldLevel) bool {
	deviceType := entities.DeviceType(fl.Field().String())
	if deviceType == entities.MOBILE || deviceType == entities.COMPUTER {
		return true
	}
	return false
}
