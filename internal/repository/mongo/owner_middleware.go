package mongo

import (
	"context"

	"github.com/pawverse/pawcare-profiles/internal/owner/domain"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"go.uber.org/zap"
)

type ownerMiddleware func(domain.IOwnerRepository) domain.IOwnerRepository

type ownerLoggingMiddleware struct {
	logger *zap.Logger
	next   domain.IOwnerRepository
}

func newOwnerLoggingMiddleware(logger *zap.Logger) ownerMiddleware {
	return func(next domain.IOwnerRepository) domain.IOwnerRepository {
		return &ownerLoggingMiddleware{logger, next}
	}
}

func (mw *ownerLoggingMiddleware) Logger(ctx context.Context) *zap.Logger {
	requestId := ctx.Value(utils.RequestIdContextKey).(string)
	return mw.logger.With(zap.String("request_id", requestId))
}

func (mw *ownerLoggingMiddleware) FindById(ctx context.Context, id domain.OwnerId) (owner *domain.Owner, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindById",
				zap.String("method", "FindById"),
				zap.String("id", string(id)),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()
	return mw.next.FindById(ctx, id)
}

func (mw *ownerLoggingMiddleware) FindByUserId(ctx context.Context, userId domain.UserId) (owner *domain.Owner, err error) {
	defer func() {
		mw.Logger(ctx).
			Info("FindByUserId",
				zap.String("method", "FindByUserId"),
				zap.String("userId", string(userId)),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()
	return mw.next.FindByUserId(ctx, userId)
}

func (mw *ownerLoggingMiddleware) Create(ctx context.Context, owner *domain.Owner) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Create",
				zap.String("method", "Create"),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()

	return mw.next.Create(ctx, owner)
}

func (mw *ownerLoggingMiddleware) Update(ctx context.Context, owner *domain.Owner) (err error) {
	defer func() {
		mw.Logger(ctx).
			Info("Update",
				zap.String("method", "Update"),
				zap.Stringer("owner", owner),
				zap.Error(err),
			)
	}()

	return mw.next.Update(ctx, owner)
}
