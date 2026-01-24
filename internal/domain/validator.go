package domain

type Validator interface {
	Validate(i any) error
	TranslateError(err error) []string
}
