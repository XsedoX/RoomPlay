package validation

import (
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"xsedox.com/main/domain/user"
)

var (
	Uni               *ut.UniversalTranslator
	ValidatorInstance *validator.Validate
	trans             ut.Translator
	whitespaceRe      = regexp.MustCompile(`\s`)
)

func Initialize() {
	ValidatorInstance = validator.New(validator.WithRequiredStructEnabled())
	
	english := en.New()
	Uni = ut.New(english, english)
	trans, _ = Uni.GetTranslator("en")

	// Register default English translations for standard tags
	if err := entranslations.RegisterDefaultTranslations(ValidatorInstance, trans); err != nil {
		panic(err)
	}
	err := ValidatorInstance.RegisterValidation("user_role_validation", validateUserRoleEnum)
	if err != nil {
		panic(err)
	}
	err = ValidatorInstance.RegisterValidation("device_type_validation", validateDeviceType)
	if err != nil {
		panic(err)
	}

	registerCustomValidation("user_role_validation", "{0} is not a valid user role", trans, validateUserRoleEnum)
	registerCustomValidation("device_type_validation", "{0} is not a valid device type", trans, validateDeviceType)
	registerCustomValidation("no_whitespace", "{0} contains spaces.", trans, validateWhitespace)

}
func registerCustomValidation(tag string, translation string, translator ut.Translator, fn validator.Func) {
	if err := ValidatorInstance.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
	if err := ValidatorInstance.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
}
func validateUserRoleEnum(fl validator.FieldLevel) bool {
	stringValue := fl.Field().String()
	return user.ParseUserRole(&stringValue) != nil
}
func validateDeviceType(fl validator.FieldLevel) bool {
	stringValue := fl.Field().String()
	return user.ParseDeviceType(&stringValue) != nil
}
func validateWhitespace(fl validator.FieldLevel) bool {
	return !whitespaceRe.MatchString(fl.Field().String())
}
func MapValidationErrors(errs validator.ValidationErrors) map[string]string {
	return errs.Translate(trans)
}
