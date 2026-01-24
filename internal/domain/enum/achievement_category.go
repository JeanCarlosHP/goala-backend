package enum

type AchievementCategory string

const (
	AchievementCategoryStreak   AchievementCategory = "streak"
	AchievementCategoryMeals    AchievementCategory = "meals"
	AchievementCategoryCalories AchievementCategory = "calories"
	AchievementCategoryProtein  AchievementCategory = "protein"
	AchievementCategoryGeneral  AchievementCategory = "general"
)

func (a AchievementCategory) IsValid() bool {
	switch a {
	case AchievementCategoryStreak, AchievementCategoryMeals,
		AchievementCategoryCalories, AchievementCategoryProtein, AchievementCategoryGeneral:
		return true
	}
	return false
}

func (a AchievementCategory) String() string {
	return string(a)
}
