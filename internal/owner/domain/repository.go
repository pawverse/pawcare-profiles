package domain

import (
	"context"
	"errors"
)

var (
	ErrOwnerNotFound       = errors.New("repository: owner not found")
	ErrOwnerAlreadyCreated = errors.New("repository: owner already created")
)

type IOwnerRepository interface {
	FindById(ctx context.Context, id OwnerId) (*Owner, error)
	FindByUserId(ctx context.Context, userId UserId) (*Owner, error)
	Create(ctx context.Context, owner *Owner) error
	Update(ctx context.Context, owner *Owner) error
}
