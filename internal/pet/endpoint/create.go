package endpoint

import (
	"context"
	"time"

	"github.com/pawverse/pawcare-profiles/internal/pet/domain"
	petservice "github.com/pawverse/pawcare-profiles/internal/pet/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Name        string    `json:"name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
	Species     string    `json:"species" validate:"required"`
	Breed       string    `json:"breed" validate:"required"`
	Weight      float64   `json:"weight" validate:"required"`
	Gender      string    `json:"gender" validate:"required,oneof=male female other"`
}

func makePetCreateEndpoint(petService petservice.IPetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return GetResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		profile := domain.NewPetProfile(req.Name, req.DateOfBirth, req.Species, req.Breed, req.Weight, domain.EGender(req.Gender))
		pet, err := petService.Create(ctx, profile)
		if err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetResponse{
			Id:          string(pet.Id),
			Name:        pet.Profile.Name,
			DateOfBirth: pet.Profile.DateOfBirth,
			Species:     pet.Profile.Species,
			Breed:       pet.Profile.Breed,
			Weight:      pet.Profile.Weight,
			Gender:      string(pet.Profile.Gender),
		}, nil
	}
}
