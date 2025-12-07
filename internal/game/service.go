package game

import (
	"context"
	"errors"
	"fmt"
	"steam_checker/internal/game/event"
	steamSchema "steam_checker/internal/shared/schema/steam"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Game, error)
	Create(ctx context.Context, input *Game) error
	Exists(ctx context.Context, appID int) (bool, error)
	GetByAppID(ctx context.Context, appID int) (Game, error)
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

func (s *Service) Create(ctx context.Context, input *CreateInput) (Game, error) {
	var (
		model Game
		err   error
	)

	if err = input.Validate(); err != nil {
		return model, err
	}

	exists, err := s.repository.Exists(ctx, input.AppID)
	if err != nil {
		return model, fmt.Errorf("erro validating if game exists: %w", err)
	}

	if exists {
		return model, errors.New("game already exists")
	}

	appDetailsResponse, err := s.steamIntegration.GetAppDetails(ctx, input.AppID)
	if err != nil {
		return model, fmt.Errorf("erro getting app details: %w", err)
	}

	playerCountResponse, err := s.steamIntegration.GetPlayerCount(ctx, input.AppID)
	if err != nil {
		return model, fmt.Errorf("erro getting player count: %w", err)
	}

	id := uuid.New()
	model = Game{
		ID:       id,
		AppID:    input.AppID,
		Name:     appDetailsResponse.Name,
		Packages: appDetailsResponse.Packages,
		Events:   s.getEvents(id, playerCountResponse, appDetailsResponse),
	}
	if err = s.repository.Create(ctx, &model); err != nil {
		return model, fmt.Errorf("error creating game: %w", err)
	}

	return model, nil
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

func (s *Service) GetByAppID(ctx context.Context, appID int) (Game, error) {
	model, err := s.repository.GetByAppID(ctx, appID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Game{}, nil
		}
		return Game{}, fmt.Errorf("error getting game by app id: %w", err)
	}

	return model, nil
}

func (s *Service) GetOrCreate(ctx context.Context, appID int) (Game, error) {
	var (
		model Game
		err   error
	)

	model, err = s.GetByAppID(ctx, appID)
	if err != nil {
		return model, err
	}

	if model.ID != uuid.Nil {
		return model, nil
	}

	model, err = s.Create(ctx, &CreateInput{
		AppID: appID,
	})
	if err != nil {
		return model, err
	}

	return model, nil
}

func (s *Service) getEvents(
	id uuid.UUID,
	playerCountResponse steamSchema.GetPlayerCountData,
	appDetailsResponse steamSchema.GetAppDetailsData,
) []event.Event {
	var events []event.Event

	playerCountEvent, err := event.New(uuid.New(), id, event.PlayerCount, playerCountResponse)
	if err == nil {
		events = append(events, *playerCountEvent)
	}

	priceEvent, err := event.New(uuid.New(), id, event.PriceUpdated, appDetailsResponse.PriceOverview)
	if err == nil {
		events = append(events, *priceEvent)
	}

	return events
}
