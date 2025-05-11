//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/pawverse/pawcare-profiles/internal/biz"
	"github.com/pawverse/pawcare-profiles/internal/conf"
	"github.com/pawverse/pawcare-profiles/internal/data"
	"github.com/pawverse/pawcare-profiles/internal/server"
	"github.com/pawverse/pawcare-profiles/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(context.Context, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
