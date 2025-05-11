package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type IOwnerRepository interface {
	Save(context.Context, *Owner) error
	Update(context.Context, *Owner) error
	FindById(context.Context, OwnerId) (*Owner, error)
	FindByUserId(context.Context, UserId) (*Owner, error)
}

type OwnerUsecase struct {
	repository IOwnerRepository
	logger     *log.Helper
}

func NewOwnerUsecase(repo IOwnerRepository, logger log.Logger) *OwnerUsecase {
	return &OwnerUsecase{
		repository: repo,
		logger:     log.NewHelper(logger),
	}
}

func (u *OwnerUsecase) Save(ctx context.Context, owner *Owner) error {
	return u.repository.Save(ctx, owner)
}

func (u *OwnerUsecase) FindById(ctx context.Context, id OwnerId) (*Owner, error) {
	return u.repository.FindById(ctx, id)
}

func (u *OwnerUsecase) FindByUserId(ctx context.Context, id UserId) (*Owner, error) {
	return u.repository.FindByUserId(ctx, id)
}
