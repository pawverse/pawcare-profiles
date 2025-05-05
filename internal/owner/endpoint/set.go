package endpoint

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/pawverse/pawcare-profiles/internal/owner/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Set struct {
	CreateEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
}

func NewSet(viper viper.Viper, ownerService service.IOwnerService) (Set, error) {
	kf, err := keyfunc.NewDefault([]string{viper.GetString(common.CertsEndpointKey)})
	if err != nil {
		return Set{}, err
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = makeCreateEndpoint(ownerService)
		createEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(createEndpoint)
	}

	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = makeGetEndpoint(ownerService)
		getEndpoint = common.NewParser(kf.Keyfunc, jwt.SigningMethodRS256, common.RegisteredClaimsFactory)(getEndpoint)
	}

	return Set{
		CreateEndpoint: createEndpoint,
		GetEndpoint:    getEndpoint,
	}, nil
}
