package analyzer

// Service defines the interface for analyzing food images.
// The implementation should return nutrition information (mocked or real).
type Service interface {
	// Analyze returns the result of a food analysis.
	Analyze() AnalysisResult
}

func NewService() Service {
	return &service{}
}
