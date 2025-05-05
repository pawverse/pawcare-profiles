package service

import (
	"github.com/pawverse/pawcare-profiles/internal/pet/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type middleware func(IPetService) IPetService

type loggingMiddleware struct {
	logger *zap.Logger
	next   IPetService
}

func newLoggingMiddleware(logger *zap.Logger) middleware {
	return func(next IPetService) IPetService {
		return &loggingMiddleware{logger, next}
	}
}

func (mw *loggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *loggingMiddleware) GetById(ctx context.Context, id domain.PetId) (pet *domain.Pet, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("GetById",
				zap.String("method", "GetById"),
				zap.String("id", string(id)),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()
	return mw.next.GetById(ctx, id)
}

func (mw *loggingMiddleware) GetAll(ctx context.Context) (pets []*domain.Pet, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("GetAll",
				zap.String("method", "GetAll"),
				zap.Int("#pets", len(pets)),
				zap.Error(err),
			)
	}()
	return mw.next.GetAll(ctx)
}

func (mw *loggingMiddleware) Create(ctx context.Context, petProfile domain.PetProfile) (pet *domain.Pet, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()

	return mw.next.Create(ctx, petProfile)
}
