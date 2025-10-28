package validation

import (
	"github.com/go-playground/validator/v10"
	"xsedox.com/main/domain/device"
	"xsedox.com/main/domain/user"
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
	_, ok := user.ParseRole(fl.Field().String())
	return ok
}
func validateDeviceType(fl validator.FieldLevel) bool {
	_, ok := device.ParseType(fl.Field().String())
	return ok
}
