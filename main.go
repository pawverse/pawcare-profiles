package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"

	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/db/mongodb"
	"github.com/pawverse/pawcare-profiles/internal/config"
	ownerendpoint "github.com/pawverse/pawcare-profiles/internal/owner/endpoint"
	ownerservice "github.com/pawverse/pawcare-profiles/internal/owner/service"
	petendpoint "github.com/pawverse/pawcare-profiles/internal/pet/endpoint"
	petservice "github.com/pawverse/pawcare-profiles/internal/pet/service"
	"github.com/pawverse/pawcare-profiles/internal/repository/mongo"
	"github.com/pawverse/pawcare-profiles/internal/transport"

	"github.com/joho/godotenv"
)

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func _main() error {
	godotenv.Load()

	viper := config.InitConfig()
	ctx := context.Background()

	http.DefaultTransport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = viper.GetBool(common.InsecureSkipVerifyKey)

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync()

	db, teardown, err := mongodb.ConnectDB(ctx, viper.GetString(common.DBConnectionStringKey), "accounts")
	defer teardown(ctx)
	if err != nil {
		return err
	}

	eventBus, err := transport.NewEventBus(viper, logger.With(zap.String("transport", "cqrs")))
	if err != nil {
		return err
	}

	ownerRepository := mongo.NewOwnerRepository(db, logger.With(zap.String("repository", "owner")))
	petRepository := mongo.NewPetRepository(db, logger.With(zap.String("repository", "pet")))

	ownerService := ownerservice.NewOwnerService(ownerRepository, logger.With(zap.String("service", "owner")))
	petService := petservice.NewPetService(petRepository, ownerService, eventBus, logger.With(zap.String("service", "pet")))

	ownerEndpoints, err := ownerendpoint.NewSet(viper, ownerService)
	if err != nil {
		return err
	}

	petEndpoints, err := petendpoint.NewSet(viper, petService)
	if err != nil {
		return err
	}

	httpHandler := transport.MakeHTTPServer(ownerEndpoints, petEndpoints, logger.With(zap.String("transport", "http")))
	httpHandler = accessControl(httpHandler)

	var g group.Group

	{
		httpAddr := fmt.Sprintf(":%s", viper.GetString(common.HTTPPortKey))
		logger := logger.With(zap.String("transport", "http"))
		httpListenAddr, err := net.Listen("tcp", httpAddr)
		if err != nil {
			return err
		}

		g.Add(func() error {
			logger.Info("Starting server", zap.String("addr", httpAddr))
			return http.Serve(httpListenAddr, httpHandler)
		}, func(err error) {
			logger.Info("Closing server", zap.NamedError("reason", err))
			httpListenAddr.Close()
		})
	}

	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		return fmt.Errorf("%s", <-c)
	}, func(err error) {
		logger.Info("Shutdown signal received", zap.NamedError("signal", err))
	})

	return g.Run()
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
