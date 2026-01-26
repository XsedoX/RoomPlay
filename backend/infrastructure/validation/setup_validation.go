package validation

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Uni               *ut.UniversalTranslator
	ValidatorInstance *validator.Validate
	trans             ut.Translator
	whitespaceRe      = regexp.MustCompile(`\s`)
	niceFieldNames    = sync.Map{}
)

func Initialize() {
	ValidatorInstance = validator.New(validator.WithRequiredStructEnabled())

	registerCustomFieldsNames()

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

	addCustomTranslation("required", "[FieldName] is required.")
	addCustomTranslation("gte", "[FieldName] must be greater than or equal to {0}.")
	addCustomTranslation("lte", "[FieldName] must be less than or equal to {0}.")
	addCustomTranslation("eqcsfield", "Values must be equal.")

	registerCustomValidation("user_role_validation", "Not a valid user role.", validateUserRoleEnum)
	registerCustomValidation("device_type_validation", "Not a valid device type.", validateDeviceType)
	registerCustomValidation("no_whitespace", "[FieldName] can not contain spaces.", validateWhitespace)
	registerCustomValidation("song_query_validation",
		"The link is not a youtube/spotify link or is incorrect.",
		validateSongQuery)
}

func addCustomTranslation(tag string, translation string) {
	if err := ValidatorInstance.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Param())
		jsonName := fe.Field()
		if niceFieldName, ok := niceFieldNames.Load(jsonName); ok {
			t = strings.ReplaceAll(t, "[FieldName]", niceFieldName.(string))
		}
		return t
	}); err != nil {
		panic(err)
	}
}

func registerCustomValidation(tag string, translation string, fn validator.Func) {
	if err := ValidatorInstance.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
	addCustomTranslation(tag, translation)
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

func validateSongQuery(fl validator.FieldLevel) bool {
	query := fl.Field().String()
	if err := ValidatorInstance.Var(query, "https_url"); err == nil {
		if !strings.Contains(query, "youtube.com") ||
			!strings.Contains(query, "youtu.be") ||
			!strings.Contains(query, "spotify.com") {
			return false
		}
	}
	return true
}

func MapValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)
	for _, e := range errs {
		errorMap[e.Field()] = e.Translate(trans)
	}
	return errorMap
}

func registerCustomFieldsNames() {
	// NOTE: used to map validation key as json tag value instead of struct field name
	// and for discovering nice names for error messages
	ValidatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			name = fld.Name
		}
		fname := fld.Tag.Get("fname")
		if fname != "" {
			niceFieldNames.Store(name, fname)
		}
		return name
	})
}
