package event

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type EventType int

const (
	PriceUpdated EventType = 0
	PlayerCount  EventType = 1
)

type Event struct {
	ID     uuid.UUID       `json:"id" db:"id"`
	GameID uuid.UUID       `json:"game_id" db:"game_id"`
	Type   EventType       `json:"type" db:"type"`
	Data   json.RawMessage `json:"data" db:"data"`
}

func New(id, gameID uuid.UUID, et EventType, data any) (*Event, error) {
	var (
		body []byte
		err  error
	)

	if data != "" {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling data: %x", err)
		}
	}

	return &Event{
		ID:     id,
		GameID: gameID,
		Type:   et,
		Data:   body,
	}, nil
}
