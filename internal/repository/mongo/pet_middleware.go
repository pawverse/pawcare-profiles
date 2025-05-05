package mongo

import (
	"context"

	ownerdomain "github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"github.com/pawverse/pawcare-profiles/internal/pet/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
)

type petMiddleware func(domain.IPetRepository) domain.IPetRepository

type petLoggingMiddleware struct {
	logger *zap.Logger
	next   domain.IPetRepository
}

func newPetLoggingMiddleware(logger *zap.Logger) petMiddleware {
	return func(next domain.IPetRepository) domain.IPetRepository {
		return &petLoggingMiddleware{logger, next}
	}
}

func (mw *petLoggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *petLoggingMiddleware) Create(ctx context.Context, pet *domain.Pet) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()

	return mw.next.Create(ctx, pet)
}

func (mw *petLoggingMiddleware) FindById(ctx context.Context, id domain.PetId) (pet *domain.Pet, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindById",
				zap.String("method", "FindById"),
				zap.String("id", string(id)),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()

	return mw.next.FindById(ctx, id)
}

func (mw *petLoggingMiddleware) FindByOwnerId(ctx context.Context, ownerId ownerdomain.OwnerId) (pets []*domain.Pet, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindByOwnerId",
				zap.String("method", "FindByOwnerId"),
				zap.String("ownerId", string(ownerId)),
				zap.Int("#pets", len(pets)),
				zap.Error(err),
			)
	}()

	return mw.next.FindByOwnerId(ctx, ownerId)
}

func (mw *petLoggingMiddleware) Update(ctx context.Context, pet *domain.Pet) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Update",
				zap.String("method", "Update"),
				zap.Stringer("pet", pet),
				zap.Error(err),
			)
	}()

	return mw.next.Update(ctx, pet)
}
