package endpoint

import (
	"context"
	"time"

	"github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"github.com/pawverse/pawcare-profiles/internal/owner/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

type GetResponse struct {
	common.EmbedError

	Id          string    `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (r GetResponse) StatusCode() int {
	switch r.Err {
	case domain.ErrOwnerNotFound:
		return 404
	case nil:
		return 200
	default:
		return 500
	}
}

var (
	_ endpoint.Failer  = (*GetResponse)(nil)
	_ http.StatusCoder = (*GetResponse)(nil)
)

func makeGetEndpoint(ownerService service.IOwnerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		owner, err := ownerService.Get(ctx)
		if err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetResponse{
			Id:          string(owner.Id),
			Name:        owner.Profile.Name,
			DateOfBirth: owner.Profile.DateOfBirth,
		}, nil
	}
}
