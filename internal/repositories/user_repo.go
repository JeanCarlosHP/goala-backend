package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "Create")
	defer span.End()

	result, err := r.db.Querier.CreateUser(ctx, db.CreateUserParams{
		ID:          pgtype.UUID{Bytes: user.ID, Valid: true},
		FirebaseUid: user.FirebaseUID,
		Email:       stringToPtr(user.Email),
		DisplayName: stringToPtr(user.DisplayName),
		PhotoUrl:    user.PhotoURL,
	})
	if err != nil {
		return err
	}

	user.CreatedAt = timePtrValue(result.CreatedAt)
	user.UpdatedAt = timePtrValue(result.UpdatedAt)
	return nil
}

func (r *UserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	tr := otel.Tracer("repositories/user_repo.go")
	ctx, span := tr.Start(ctx, "GetByFirebaseUID")
	defer span.End()

	result, err := r.db.Querier.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:                   result.ID.Bytes,
		FirebaseUID:          result.FirebaseUid,
		Email:                stringPtrValue(result.Email),
		DisplayName:          stringPtrValue(result.DisplayName),
		PhotoURL:             result.PhotoUrl,
		Weight:               intPtrToInt32Ptr(result.Weight),
		Height:               intPtrToInt32Ptr(result.Height),
		Age:                  intPtrToInt32Ptr(result.Age),
		Gender:               result.Gender,
		ActivityLevel:        result.ActivityLevel,
		Language:             stringPtrValue(result.Language),
		NotificationsEnabled: boolPtrValue(result.NotificationsEnabled),
		CreatedAt:            timePtrValue(result.CreatedAt),
		UpdatedAt:            timePtrValue(result.UpdatedAt),
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "GetByID")
	defer span.End()

	result, err := r.db.Querier.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:                   result.ID.Bytes,
		FirebaseUID:          result.FirebaseUid,
		Email:                stringPtrValue(result.Email),
		DisplayName:          stringPtrValue(result.DisplayName),
		PhotoURL:             result.PhotoUrl,
		Weight:               intPtrToInt32Ptr(result.Weight),
		Height:               intPtrToInt32Ptr(result.Height),
		Age:                  intPtrToInt32Ptr(result.Age),
		Gender:               result.Gender,
		ActivityLevel:        result.ActivityLevel,
		Language:             stringPtrValue(result.Language),
		NotificationsEnabled: boolPtrValue(result.NotificationsEnabled),
		CreatedAt:            timePtrValue(result.CreatedAt),
		UpdatedAt:            timePtrValue(result.UpdatedAt),
	}, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "Update")
	defer span.End()

	result, err := r.db.Querier.UpdateUser(ctx, db.UpdateUserParams{
		ID:          pgtype.UUID{Bytes: user.ID, Valid: true},
		Email:       stringToPtr(user.Email),
		DisplayName: stringToPtr(user.DisplayName),
		PhotoUrl:    user.PhotoURL,
	})
	if err != nil {
		return err
	}

	user.UpdatedAt = timePtrValue(result.UpdatedAt)
	return nil
}

func (r *UserRepository) ExistsByFirebaseUID(ctx context.Context, firebaseUID string) (bool, error) {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "ExistsByFirebaseUID")
	defer span.End()

	return r.db.Querier.ExistsUserByFirebaseUID(ctx, firebaseUID)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "UpdateProfile")
	defer span.End()

	return r.db.Querier.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:                   pgtype.UUID{Bytes: user.ID, Valid: true},
		DisplayName:          stringToPtr(user.DisplayName),
		Email:                stringToPtr(user.Email),
		PhotoUrl:             user.PhotoURL,
		Weight:               int32PtrToIntPtr(user.Weight),
		Height:               int32PtrToIntPtr(user.Height),
		Age:                  int32PtrToIntPtr(user.Age),
		Gender:               user.Gender,
		ActivityLevel:        user.ActivityLevel,
		Language:             stringToPtr(user.Language),
		NotificationsEnabled: boolToPtr(user.NotificationsEnabled),
	})
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, userID uuid.UUID, photoURL *string) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "UpdateAvatar")
	defer span.End()

	return r.db.Querier.UpdateUserAvatar(ctx, db.UpdateUserAvatarParams{
		ID:       pgtype.UUID{Bytes: userID, Valid: true},
		PhotoUrl: photoURL,
	})
}

func (r *UserRepository) UpdateDisplayName(ctx context.Context, userID uuid.UUID, displayName *string) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "UpdateDisplayName")
	defer span.End()

	return r.db.Querier.UpdateUserDisplayName(ctx, db.UpdateUserDisplayNameParams{
		ID:          pgtype.UUID{Bytes: userID, Valid: true},
		DisplayName: displayName,
	})
}

func (r *UserRepository) UpdateNotifications(ctx context.Context, userID uuid.UUID, notificationsEnabled *bool) error {
	tr := otel.Tracer("services/user_repo.go")
	ctx, span := tr.Start(ctx, "UpdateNotifications")
	defer span.End()

	return r.db.Querier.UpdateUserNotifications(ctx, db.UpdateUserNotificationsParams{
		ID:                   pgtype.UUID{Bytes: userID, Valid: true},
		NotificationsEnabled: notificationsEnabled,
	})
}
