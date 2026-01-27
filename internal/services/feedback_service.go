package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"
)

type FeedbackService struct {
	feedbackRepo *repositories.FeedbackRepository
	logger       domain.Logger
}

func NewFeedbackService(feedbackRepo *repositories.FeedbackRepository, logger domain.Logger) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
		logger:       logger,
	}
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, userID uuid.UUID, req *domain.CreateFeedbackRequest) error {
	tr := otel.Tracer("services/feedback_service.go")
	ctx, span := tr.Start(ctx, "CreateFeedback")
	defer span.End()

	feedback, err := s.feedbackRepo.Create(ctx, userID, req)
	if err != nil {
		s.logger.Error("Failed to create feedback", "user_id", userID.String(), "error", err)
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	s.logger.Info("Feedback created successfully", "feedback_id", feedback.ID, "user_id", userID.String(), "type", string(feedback.Type), "title", feedback.Title)

	return nil
}

func (s *FeedbackService) GetFeedback(ctx context.Context, feedbackID uuid.UUID) (*domain.Feedback, error) {
	tr := otel.Tracer("services/feedback_service.go")
	ctx, span := tr.Start(ctx, "GetFeedback")
	defer span.End()

	feedback, err := s.feedbackRepo.GetByID(ctx, feedbackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return feedback, nil
}

func (s *FeedbackService) ListFeedback(ctx context.Context, limit, offset int32) ([]domain.Feedback, error) {
	tr := otel.Tracer("services/feedback_service.go")
	ctx, span := tr.Start(ctx, "ListFeedback")
	defer span.End()

	feedbacks, err := s.feedbackRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedback: %w", err)
	}

	return feedbacks, nil
}

func (s *FeedbackService) GetUserFeedback(ctx context.Context, userID uuid.UUID) ([]domain.Feedback, error) {
	tr := otel.Tracer("services/feedback_service.go")
	ctx, span := tr.Start(ctx, "GetUserFeedback")
	defer span.End()

	feedbacks, err := s.feedbackRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feedback: %w", err)
	}

	return feedbacks, nil
}
