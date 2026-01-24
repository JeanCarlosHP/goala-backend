package enum

type ActivityLevel string

const (
	ActivityLevelSedentary  ActivityLevel = "sedentary"
	ActivityLevelLight      ActivityLevel = "light"
	ActivityLevelModerate   ActivityLevel = "moderate"
	ActivityLevelActive     ActivityLevel = "active"
	ActivityLevelVeryActive ActivityLevel = "very_active"
)

func (a ActivityLevel) IsValid() bool {
	switch a {
	case ActivityLevelSedentary, ActivityLevelLight, ActivityLevelModerate,
		ActivityLevelActive, ActivityLevelVeryActive:
		return true
	}
	return false
}

func (a ActivityLevel) String() string {
	return string(a)
}
