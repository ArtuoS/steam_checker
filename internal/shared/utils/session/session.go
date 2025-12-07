package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	idStr := ctx.Value("user_id").(string)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
