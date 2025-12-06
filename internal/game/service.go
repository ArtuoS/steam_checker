package game

import (
	"context"
	"errors"
	"fmt"
	"steam_checker/internal/game/event"
	steamSchema "steam_checker/internal/shared/schema/steam"

	"github.com/google/uuid"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Game, error)
	Create(ctx context.Context, input *Game) error
	Exists(ctx context.Context, appID int) (bool, error)
}

type EventRepository interface {
	GetByGameID(ctx context.Context, gameID uuid.UUID) ([]event.Event, error)
}

type SteamIntegration interface {
	GetPlayerCount(ctx context.Context, appID int) (steamSchema.GetPlayerCountData, error)
	GetAppDetails(ctx context.Context, appID int) (steamSchema.GetAppDetailsData, error)
}

type Service struct {
	repository       Repository
	eventRepository  EventRepository
	steamIntegration SteamIntegration
}

func NewService(repository Repository, eventRepository EventRepository, steamIntegration SteamIntegration) *Service {
	return &Service{
		repository:       repository,
		eventRepository:  eventRepository,
		steamIntegration: steamIntegration,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	exists, err := s.repository.Exists(ctx, input.AppID)
	if err != nil {
		return fmt.Errorf("erro validating if game exists: %w", err)
	}

	if exists {
		return errors.New("game already exists")
	}

	appDetailsResponse, err := s.steamIntegration.GetAppDetails(ctx, input.AppID)
	if err != nil {
		return fmt.Errorf("erro getting app details: %w", err)
	}

	playerCountResponse, err := s.steamIntegration.GetPlayerCount(ctx, input.AppID)
	if err != nil {
		return fmt.Errorf("erro getting player count: %w", err)
	}

	id := uuid.New()
	evt, err := event.New(uuid.New(), id, event.PlayerCount, playerCountResponse)
	if err != nil {
		return fmt.Errorf("erro creating event: %w", err)
	}

	err = s.repository.Create(ctx, &Game{
		ID:       id,
		AppID:    input.AppID,
		Name:     appDetailsResponse.Name,
		Packages: appDetailsResponse.Packages,
		Events: []event.Event{
			*evt,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating game: %w", err)
	}

	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]Game, error) {
	models, err := s.repository.GetAll(ctx)
	if err != nil {
		return models, fmt.Errorf("error getting all games: %w", err)
	}

	for _, mdl := range models {
		events, err := s.eventRepository.GetByGameID(ctx, mdl.ID)
		if err != nil {
			return models, err
		}

		mdl.Events = append(mdl.Events, events...)
	}

	return models, nil
}
