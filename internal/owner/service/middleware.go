package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"go.uber.org/zap"
)

type middleware func(IOwnerService) IOwnerService

type loggingMiddleware struct {
	logger *zap.Logger
	next   IOwnerService
}

func newLoggingMiddleware(logger *zap.Logger) middleware {
	return func(next IOwnerService) IOwnerService {
		return &loggingMiddleware{logger, next}
	}
}

func (mw *loggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *loggingMiddleware) Get(ctx context.Context) (owner *domain.Owner, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Get",
				zap.String("method", "Get"),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()

	return mw.next.Get(ctx)
}

func (mw *loggingMiddleware) Create(ctx context.Context, ownerProfile domain.OwnerProfile) (owner *domain.Owner, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.Stringer("ownerProfile", ownerProfile),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()

	return mw.next.Create(ctx, ownerProfile)
}

type validationMiddleware struct {
	next IOwnerService
}

func newValidationMiddleware() middleware {
	return func(next IOwnerService) IOwnerService {
		return &validationMiddleware{next}
	}
}

func (mw *validationMiddleware) Get(ctx context.Context) (*domain.Owner, error) {
	return mw.next.Get(ctx)
}

func (mw *validationMiddleware) Create(ctx context.Context, ownerProfile domain.OwnerProfile) (*domain.Owner, error) {
	if err := validator.New().Struct(ownerProfile); err != nil {
		return nil, err
	}

	if time.Now().Before(ownerProfile.DateOfBirth) || time.Now().AddDate(-150, 0, 0).After(ownerProfile.DateOfBirth) {
		return nil, ErrInvalidDate
	}

	return mw.next.Create(ctx, ownerProfile)
}
