package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type validator struct {
	validate  *v10.Validate
	translate ut.Translator
}

func NewValidator() domain.Validator {
	en := en.New()
	uni := ut.New(en, en)

	translate, _ := uni.GetTranslator("en")

	v := v10.New()
	err := en_translations.RegisterDefaultTranslations(v, translate)
	if err != nil {
		panic(err)
	}

	return &validator{
		validate:  v,
		translate: translate,
	}
}

func (v *validator) Validate(i any) error {
	return v.validate.Struct(i)
}

func (v *validator) TranslateError(err error) []string {
	if err == nil {
		return nil
	}

	errs := err.(v10.ValidationErrors)
	translations := errs.Translate(v.translate)

	var errorMessages []string
	for _, translation := range translations {
		formattedMessage := formatTranslation(translation)
		errorMessages = append(errorMessages, formattedMessage)
	}

	return errorMessages
}

func formatTranslation(translation string) string {
	translation = fmt.Sprintf("%s%s", strings.ToLower(string(translation[0])), translation[1:])
	translation = strings.ReplaceAll(translation, "ID", "Id")

	return translation
}
