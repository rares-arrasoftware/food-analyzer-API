package analyzer

type service struct{}

func (s *service) Analyze() AnalysisResult {
	// Stubbed response
	return AnalysisResult{
		ItemName:   "Grilled Chicken Breast",
		Calories:   231,
		Protein:    43.5,
		Carbs:      0,
		Fat:        5.0,
		Confidence: 0.95,
	}
}
