package enum

type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeSnack     MealType = "snack"
)

func (m MealType) IsValid() bool {
	switch m {
	case MealTypeBreakfast, MealTypeLunch, MealTypeDinner, MealTypeSnack:
		return true
	}
	return false
}

func (m MealType) String() string {
	return string(m)
}
