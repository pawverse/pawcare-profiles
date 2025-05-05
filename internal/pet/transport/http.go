package transport

import (
	"github.com/pawverse/pawcare-profiles/internal/pet/endpoint"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/utils"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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

	authenticatedOpts := append(options, kithttp.ServerBefore(kitjwt.HTTPToContext()))

	createHandler := kithttp.NewServer(
		endpoints.CreateEndpoint,
		common.DecodeJSONRequest[endpoint.CreateRequest],
		common.EncodeJSONResponse,
		authenticatedOpts...,
	)

	getAllHandler := kithttp.NewServer(
		endpoints.GetAllEndpoint,
		common.DecodeNoBodyRequest,
		common.EncodeJSONResponse,
		authenticatedOpts...,
	)

	getByIdHandler := kithttp.NewServer(
		endpoints.GetByIdEndpoint,
		common.DecodePathParameters[endpoint.GetByIdRequest],
		common.EncodeJSONResponse,
		authenticatedOpts...,
	)

	router := r.PathPrefix("/pets").Subrouter()
	router.Methods("POST").Path("").Handler(createHandler)
	router.Methods("GET").Path("").Handler(getAllHandler)
	router.Methods("GET").Path("/{id}").Handler(getByIdHandler)
}
