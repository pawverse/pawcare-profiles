package endpoint

import (
	"github.com/MicahParks/keyfunc/v3"
	petservice "github.com/pawverse/pawcare-profiles/internal/pet/service"
	"github.com/pawverse/pawcare-core/pkg/common"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Set struct {
	CreateEndpoint  endpoint.Endpoint
	GetByIdEndpoint endpoint.Endpoint
	GetAllEndpoint  endpoint.Endpoint
}

func NewSet(viper viper.Viper, petService petservice.IPetService) (Set, error) {
	kf, err := keyfunc.NewDefault([]string{viper.GetString(common.CertsEndpointKey)})
	if err != nil {
		return Set{}, err
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = makePetCreateEndpoint(petService)
		createEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(createEndpoint)
	}

	var getByIdEndpoint endpoint.Endpoint
	{
		getByIdEndpoint = makeGetByIdEndpoint(petService)
		getByIdEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(getByIdEndpoint)
	}

	var getAllEndpoint endpoint.Endpoint
	{
		getAllEndpoint = makeGetAllEndpoint(petService)
		getAllEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(getAllEndpoint)
	}

	return Set{
		CreateEndpoint:  createEndpoint,
		GetByIdEndpoint: getByIdEndpoint,
		GetAllEndpoint:  getAllEndpoint,
	}, nil
}
