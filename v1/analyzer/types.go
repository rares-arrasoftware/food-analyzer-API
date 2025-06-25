package analyzer

type AnalysisResult struct {
	ItemName   string  `json:"itemName"`
	Calories   int     `json:"calories"`
	Protein    float64 `json:"protein"`
	Carbs      float64 `json:"carbs"`
	Fat        float64 `json:"fat"`
	Confidence float64 `json:"confidence"`
}
