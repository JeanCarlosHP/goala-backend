package enum

type Language string

const (
	LanguageEnUS Language = "en-US"
	LanguagePtBR Language = "pt-BR"
)

func (l Language) IsValid() bool {
	switch l {
	case LanguageEnUS, LanguagePtBR:
		return true
	}
	return false
}

func (l Language) String() string {
	return string(l)
}
