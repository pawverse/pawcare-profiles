package transport

import (
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"github.com/pawverse/pawcare-profiles/internal/owner/endpoint"
	"go.uber.org/zap"
)

func RegisterHTTPRoutes(r *mux.Router, endpoints endpoint.Set, logger *zap.Logger) {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(common.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(utils.RequestIdHTTPToContext()),
	}
	options = append(options, common.HTTPLoggingServerOptions(logger)...)

	createHandler := kithttp.NewServer(
		endpoints.CreateEndpoint,
		common.DecodeJSONRequest[endpoint.CreateRequest],
		common.EncodeJSONResponse,
		options...,
	)

	getHandler := kithttp.NewServer(
		endpoints.GetEndpoint,
		common.DecodeNoBodyRequest,
		common.EncodeJSONResponse,
		options...,
	)

	router := r.PathPrefix("/owners").Subrouter()
	router.Methods("POST").Path("").Handler(createHandler)
	router.Methods("GET").Path("").Handler(getHandler)
}
