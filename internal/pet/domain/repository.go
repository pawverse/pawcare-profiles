package domain

import (
	"context"
	"errors"

	"github.com/pawverse/pawcare-profiles/internal/owner/domain"
)

var ErrPetNotFound = errors.New("repository: pet not found")

type IPetRepository interface {
	FindById(ctx context.Context, id PetId) (*Pet, error)
	FindByOwnerId(ctx context.Context, ownerId domain.OwnerId) ([]*Pet, error)
	Create(ctx context.Context, pet *Pet) error
	Update(ctx context.Context, pet *Pet) error
}
