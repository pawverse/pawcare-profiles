package data

import (
	"context"

	"github.com/pawverse/pawcare-profiles/internal/biz"
)

type ownerRepository struct {
	data *Data
}

func NewOwnerRepo(data *Data) biz.IOwnerRepository {
	return &ownerRepository{data}
}

// FindById implements biz.IOwnerRepository.
func (o *ownerRepository) FindById(context.Context, biz.OwnerId) (*biz.Owner, error) {
	panic("unimplemented")
}

// FindByUserId implements biz.IOwnerRepository.
func (o *ownerRepository) FindByUserId(context.Context, biz.UserId) (*biz.Owner, error) {
	panic("unimplemented")
}

// Save implements biz.IOwnerRepository.
func (o *ownerRepository) Save(context.Context, *biz.Owner) error {
	return nil
}

// Update implements biz.IOwnerRepository.
func (o *ownerRepository) Update(context.Context, *biz.Owner) error {
	panic("unimplemented")
}
