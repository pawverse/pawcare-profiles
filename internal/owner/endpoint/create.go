package endpoint

import (
	"context"
	"time"

	"github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"github.com/pawverse/pawcare-profiles/internal/owner/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"

	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Name        string    `json:"name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}

type CreateResponse struct {
	common.EmbedError

	Id   string `json:"id"`
	Name string `json:"name"`
}

func (r CreateResponse) StatusCode() int {
	switch r.Err {
	case domain.ErrOwnerAlreadyCreated:
		return 409
	case nil:
		return 201
	default:
		return 500
	}
}

var (
	_ endpoint.Failer  = (*CreateResponse)(nil)
	_ http.StatusCoder = (*CreateResponse)(nil)
)

func makeCreateEndpoint(ownerService service.IOwnerService) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return CreateResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		err := validator.New().Struct(req)
		if err != nil {
			return CreateResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		profile := domain.NewOwnerProfile(req.Name, req.DateOfBirth)
		owner, err := ownerService.Create(context, profile)
		if err != nil {
			return CreateResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return CreateResponse{
			Id:   string(owner.Id),
			Name: owner.Profile.Name,
		}, nil
	}
}
