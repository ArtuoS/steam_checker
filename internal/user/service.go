package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"steam_checker/internal/game"
	"steam_checker/internal/shared/utils/password"
	"steam_checker/internal/shared/utils/session"
	"steam_checker/internal/user/user_game"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, input *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type GameService interface {
	GetOrCreate(ctx context.Context, appID int) (game.Game, error)
}

type UserGameRepository interface {
	Create(ctx context.Context, userGame *user_game.UserGame) error
}

type Service struct {
	repository         Repository
	gameService        GameService
	userGameRepository UserGameRepository
}

func NewService(repository Repository, gameService GameService, userGameRepository UserGameRepository) *Service {
	return &Service{
		repository:         repository,
		gameService:        gameService,
		userGameRepository: userGameRepository,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	exists, err := s.repository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("error validating if user exists: %w", err)
	}

	if exists {
		return errors.New("user with this email already exists")
	}

	pswd, err := password.Encrypt(input.Password)
	if err != nil {
		return fmt.Errorf("error encrypting password: %w", err)
	}

	id := uuid.New()
	err = s.repository.Create(ctx, &User{
		ID:       id,
		Name:     input.Name,
		Email:    input.Email,
		Password: pswd,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (s *Service) Authenticate(ctx context.Context, input *AuthenticateInput) (AuthenticateOutput, error) {
	var out AuthenticateOutput

	if err := input.Validate(); err != nil {
		return out, err
	}

	user, err := s.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		return out, fmt.Errorf("error getting users: %w", err)
	}

	ok, err := password.AreEqual(input.Password, user.Password)
	if err != nil {
		return out, fmt.Errorf("error comparing passwords: %w", err)
	}
	if !ok {
		return out, errors.New("invalid credentials")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"sub":     user.ID.String(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return out, fmt.Errorf("error signing token: %w", err)
	}

	out.Token = tokenString
	return out, nil
}

func (s *Service) Track(ctx context.Context, input *TrackInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	userID, err := session.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("error getting authenticated user: %w", err)
	}

	existing, err := s.gameService.GetOrCreate(ctx, input.AppID)
	if err != nil {
		return fmt.Errorf("error getting game by app id while tracking: %w", err)
	}

	if err := s.userGameRepository.Create(ctx, &user_game.UserGame{
		ID:     uuid.New(),
		UserID: userID,
		GameID: existing.ID,
	}); err != nil {
		return fmt.Errorf("error creating user game connection: %w", err)
	}

	return nil
}
