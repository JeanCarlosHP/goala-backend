package enum

type AIFeature string

const (
	FeatureFoodRecognition AIFeature = "food_recognition"
	FeatureMealAnalysis    AIFeature = "meal_analysis"
	FeatureNutritionAdvice AIFeature = "nutrition_advice"
)

func (f AIFeature) String() string {
	return string(f)
}

func (f AIFeature) IsValid() bool {
	switch f {
	case FeatureFoodRecognition, FeatureMealAnalysis, FeatureNutritionAdvice:
		return true
	}
	return false
}
