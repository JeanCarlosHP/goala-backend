package domain

type ProgressUpdate struct {
	Stage      string `json:"stage"`
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
}

type FoodRecognitionRequest struct {
	URI          string `form:"uri" validate:"required"`
	Name         string `form:"name" validate:"required"`
	Type         string `form:"type" validate:"required"`
	MealLocation string `form:"mealLocation" validate:"required"`
}

type RecognizedFoodItem struct {
	Name       string  `json:"name" validate:"required"`
	Calories   int32   `json:"calories" validate:"gte=0,lte=5000"`
	Protein    int32   `json:"protein" validate:"gte=0,lte=500"`
	Carbs      int32   `json:"carbs" validate:"gte=0,lte=500"`
	Fat        int32   `json:"fat" validate:"gte=0,lte=500"`
	Quantity   int32   `json:"quantity" validate:"gte=1,lte=10000"`
	Unit       string  `json:"unit" validate:"required"`
	Confidence float64 `json:"confidence" validate:"gte=0,lte=1"`
}

type FoodRecognitionResponse struct {
	FoodItems      []RecognizedFoodItem `json:"foodItems" validate:"required,dive"`
	ProcessingTime int32                `json:"processingTime" validate:"gte=0"`
}

type FoodBarcodeResponse struct {
	Barcode     string  `json:"barcode" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Brand       *string `json:"brand"`
	Calories    int32   `json:"calories" validate:"required,gte=0,lte=5000"`
	Protein     int32   `json:"protein" validate:"required,gte=0,lte=5000"`
	Carbs       int32   `json:"carbs" validate:"required,gte=0,lte=5000"`
	Fat         int32   `json:"fat" validate:"required,gte=0,lte=5000"`
	ServingSize *int32  `json:"servingSize" validate:"omitempty,gte=1,lte=5000"`
	ServingUnit *string `json:"servingUnit"`
	Source      *string `json:"source"`
}

type EstimateQuantityRequest struct {
	URI                  string  `form:"uri" validate:"required"`
	Name                 string  `form:"name" validate:"required"`
	Type                 string  `form:"type" validate:"required"`
	MealLocation         string  `form:"mealLocation" validate:"required"`
	ReferenceServingSize *string `form:"referenceServingSize"`
	ReferenceServingUnit *string `form:"referenceServingUnit"`
}

type EstimateQuantityResponse struct {
	EstimatedQuantity int32   `json:"estimatedQuantity" validate:"required,gte=1,lte=500"`
	Unit              string  `json:"unit" validate:"required,oneof=g ml serving"`
	Confidence        float64 `json:"confidence" validate:"required,gte=0,lte=1"`
	Reasoning         *string `json:"reasoning"`
}

type ProcessingProgress struct {
	Stage      string                   `json:"stage"`
	Percentage int                      `json:"percentage"`
	Message    string                   `json:"message"`
	Data       *FoodRecognitionResponse `json:"data,omitempty"`
	Error      *string                  `json:"error,omitempty"`
}
