package service

import (
	"context"
	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(GatherService) GatherService

type loggingMiddleware struct {
	logger log.Logger
	next   GatherService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a GatherService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next GatherService) GatherService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) GatherSalesVolumes(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	defer func() {
		l.logger.Log("method", "GatherSalesVolumes", "targetID", targetID, "startTime", startTime, "endTime", endTime, "num", num, "err", err)
	}()
	return l.next.GatherSalesVolumes(ctx, targetID, startTime, endTime)
}
func (l loggingMiddleware) GatherCommentsNumber(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	defer func() {
		l.logger.Log("method", "GatherCommentsNumber", "targetID", targetID, "startTime", startTime, "endTime", endTime, "num", num, "err", err)
	}()
	return l.next.GatherCommentsNumber(ctx, targetID, startTime, endTime)
}
func (l loggingMiddleware) GatherViews(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	defer func() {
		l.logger.Log("method", "GatherViews", "targetID", targetID, "startTime", startTime, "endTime", endTime, "num", num, "err", err)
	}()
	return l.next.GatherViews(ctx, targetID, startTime, endTime)
}
