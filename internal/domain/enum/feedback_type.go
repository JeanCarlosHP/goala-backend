package enum

type FeedbackType string

const (
	FeedbackTypeProblem     FeedbackType = "problem"
	FeedbackTypeImprovement FeedbackType = "improvement"
)

func (f FeedbackType) IsValid() bool {
	switch f {
	case FeedbackTypeProblem, FeedbackTypeImprovement:
		return true
	}
	return false
}

func (f FeedbackType) String() string {
	return string(f)
}
