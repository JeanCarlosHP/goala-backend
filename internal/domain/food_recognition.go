package domain

type ProgressUpdate struct {
	Stage      string `json:"stage"`
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
}

type FoodRecognitionRequest struct {
	ImagePath    string `json:"imagePath" validate:"required,startswith=/"`
	Name         string `json:"name" validate:"required"`
	Type         string `json:"type" validate:"required"`
	MealLocation string `json:"mealLocation" validate:"required"`
}

type RecognizedFoodItem struct {
	Name       string  `json:"name" validate:"required"`
	Calories   float64 `json:"calories" validate:"gte=0,lte=5000"`
	Protein    float64 `json:"protein" validate:"gte=0,lte=500"`
	Carbs      float64 `json:"carbs" validate:"gte=0,lte=500"`
	Fat        float64 `json:"fat" validate:"gte=0,lte=500"`
	Quantity   float64 `json:"quantity" validate:"gte=1,lte=10000"`
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
	ImagePath            string  `json:"imagePath" validate:"required,startswith=/"`
	Name                 string  `json:"name" validate:"required"`
	Type                 string  `json:"type" validate:"required"`
	MealLocation         string  `json:"mealLocation" validate:"required"`
	ReferenceServingSize *string `json:"referenceServingSize"`
	ReferenceServingUnit *string `json:"referenceServingUnit"`
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

type FoodImageUploadRequest struct {
	ContentType string `json:"contentType" validate:"required,oneof=image/jpeg image/jpg image/png image/webp"`
	FileSize    int64  `json:"fileSize" validate:"required,gt=0,lte=5242880"`
}

type FoodImageUploadResponse struct {
	UploadURL string `json:"uploadUrl"`
	ImagePath string `json:"imagePath"`
	ExpiresIn int    `json:"expiresIn"`
}
