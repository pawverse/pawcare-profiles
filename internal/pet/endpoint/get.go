package endpoint

import (
	"context"
	"errors"
	"time"

	"github.com/pawverse/pawcare-profiles/internal/pet/domain"
	"github.com/pawverse/pawcare-profiles/internal/pet/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
)

type GetManyResponse struct {
	common.EmbedError

	Pets []GetResponse `json:"pets"`
}

func (r GetManyResponse) StatusCode() int {
	return 200
}

var _ http.StatusCoder = (*GetManyResponse)(nil)

type GetByIdRequest struct {
	Id string `json:"id" validate:"required,mongodb"`
}

type GetResponse struct {
	common.EmbedError

	Id          string    `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Species     string    `json:"species"`
	Breed       string    `json:"breed"`
	Weight      float64   `json:"weight"`
	Gender      string    `json:"gender"`
}

func (r GetResponse) StatusCode() int {
	var validationError validator.ValidationErrors
	if errors.As(r.Err, &validationError) {
		return 400
	}

	switch r.Err {
	case domain.ErrPetNotFound:
		return 404
	case nil:
		return 200
	default:
		return 500
	}
}

var _ http.StatusCoder = (*GetResponse)(nil)

func makeGetAllEndpoint(petService service.IPetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		pets, err := petService.GetAll(ctx)
		if err != nil {
			return GetManyResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetManyResponse{
			Pets: utils.Map(pets, func(pet *domain.Pet) GetResponse {
				return GetResponse{
					Id:          string(pet.Id),
					Name:        pet.Profile.Name,
					Species:     pet.Profile.Species,
					Breed:       pet.Profile.Breed,
					DateOfBirth: pet.Profile.DateOfBirth,
					Weight:      pet.Profile.Weight,
					Gender:      string(pet.Profile.Gender),
				}
			}),
		}, nil
	}
}

func makeGetByIdEndpoint(petService service.IPetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetByIdRequest)
		if !ok {
			return GetResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		pet, err := petService.GetById(ctx, domain.PetId(req.Id))
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
