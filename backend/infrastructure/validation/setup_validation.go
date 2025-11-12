package validation

import (
	"reflect"
	"regexp"
	"strings"

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

	ValidatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

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

	overrideDefaultTranslation("required", "This field is required")
	overrideDefaultTranslation("gte", "Value must be greater than or equal to {0}")
	overrideDefaultTranslation("lte", "Value must be less than or equal to {0}")
	overrideDefaultTranslation("eqcsfield", "Values must be equal")

	registerCustomValidation("user_role_validation", "Not a valid user role", trans, validateUserRoleEnum)
	registerCustomValidation("device_type_validation", "Not a valid device type", trans, validateDeviceType)
	registerCustomValidation("no_whitespace", "It can not contain spaces.", trans, validateWhitespace)

}
func overrideDefaultTranslation(tag string, translation string) {
	if err := ValidatorInstance.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Param())
		return t
	}); err != nil {
		panic(err)
	}
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
	errorMap := make(map[string]string)
	for _, e := range errs {
		errorMap[e.Field()] = e.Translate(trans)
	}
	return errorMap
}
