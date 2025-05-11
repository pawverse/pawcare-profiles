package service

import (
	"context"

	v1 "github.com/pawverse/pawcare-profiles/api/owners/v1"
	"github.com/pawverse/pawcare-profiles/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OwnerService struct {
	v1.UnimplementedOwnerServiceServer

	uc *biz.OwnerUsecase
}

func NewOwnerService(uc *biz.OwnerUsecase) *OwnerService {
	return &OwnerService{uc: uc}
}

func (s *OwnerService) CreateOwner(ctx context.Context, req *v1.CreateOwnerRequest) (*v1.CreateOwnerResponse, error) {
	owner := &biz.Owner{
		Name:        req.Name,
		UserId:      biz.UserId(""),
		DateOfBirth: req.DateOfBirth.AsTime(),
	}
	if err := s.uc.Save(ctx, owner); err != nil {
		return nil, err
	}

	return &v1.CreateOwnerResponse{
		Id:          string(owner.Id),
		Name:        owner.Name,
		DateOfBirth: timestamppb.New(owner.DateOfBirth),
	}, nil
}
