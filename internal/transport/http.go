package transport

import (
	"net/http"

	ownerendpoint "github.com/pawverse/pawcare-profiles/internal/owner/endpoint"
	ownertransport "github.com/pawverse/pawcare-profiles/internal/owner/transport"
	petendpoint "github.com/pawverse/pawcare-profiles/internal/pet/endpoint"
	pettransport "github.com/pawverse/pawcare-profiles/internal/pet/transport"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func MakeHTTPServer(ownerEndpoints ownerendpoint.Set, petEndpoints petendpoint.Set, logger *zap.Logger) http.Handler {
	router := mux.NewRouter()
	apiGroup := router.PathPrefix("/api/v1").Subrouter()

	ownertransport.RegisterHTTPRoutes(apiGroup, ownerEndpoints, logger)
	pettransport.RegisterHTTPRoutes(apiGroup, petEndpoints, logger)

	return router
}
