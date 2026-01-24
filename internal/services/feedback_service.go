package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"github.com/rs/zerolog/log"
)

type FeedbackService struct {
	feedbackRepo *repositories.FeedbackRepository
}

func NewFeedbackService(feedbackRepo *repositories.FeedbackRepository) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, userID uuid.UUID, req *domain.CreateFeedbackRequest) error {
	feedback, err := s.feedbackRepo.Create(ctx, userID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("type", req.Type).
			Msg("Failed to create feedback")
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	log.Info().
		Str("feedback_id", feedback.ID).
		Str("user_id", userID.String()).
		Str("type", string(feedback.Type)).
		Str("title", feedback.Title).
		Msg("Feedback created successfully")

	return nil
}

func (s *FeedbackService) GetFeedback(ctx context.Context, feedbackID uuid.UUID) (*domain.Feedback, error) {
	feedback, err := s.feedbackRepo.GetByID(ctx, feedbackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return feedback, nil
}

func (s *FeedbackService) ListFeedback(ctx context.Context, limit, offset int32) ([]domain.Feedback, error) {
	feedbacks, err := s.feedbackRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedback: %w", err)
	}

	return feedbacks, nil
}

func (s *FeedbackService) GetUserFeedback(ctx context.Context, userID uuid.UUID) ([]domain.Feedback, error) {
	feedbacks, err := s.feedbackRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feedback: %w", err)
	}

	return feedbacks, nil
}
