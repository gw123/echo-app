package service

import "context"

// GatherService describes the service.
type GatherService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	GatherSalesVolumes(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error)
	GatherCommentsNumber(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error)
	GatherViews(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error)
}

type basicGatherService struct{}

func (b *basicGatherService) GatherSalesVolumes(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	// TODO implement the business logic of GatherSalesVolumes
	return num, err
}
func (b *basicGatherService) GatherCommentsNumber(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	// TODO implement the business logic of GatherCommentsNumber
	return num, err
}
func (b *basicGatherService) GatherViews(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	// TODO implement the business logic of GatherViews
	return num, err
}

// NewBasicGatherService returns a naive, stateless implementation of GatherService.
func NewBasicGatherService() GatherService {
	return &basicGatherService{}
}

// New returns a GatherService with all of the expected middleware wired in.
func New(middleware []Middleware) GatherService {
	var svc GatherService = NewBasicGatherService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
