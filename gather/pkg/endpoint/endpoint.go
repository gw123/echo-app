package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	service "github.com/gw123/echo-app/gather/pkg/service"
)

// GatherSalesVolumesRequest collects the request parameters for the GatherSalesVolumes method.
type GatherSalesVolumesRequest struct {
	TargetID  int64  `json:"target_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// GatherSalesVolumesResponse collects the response parameters for the GatherSalesVolumes method.
type GatherSalesVolumesResponse struct {
	Num int64 `json:"num"`
	Err error `json:"err"`
}

// MakeGatherSalesVolumesEndpoint returns an endpoint that invokes GatherSalesVolumes on the service.
func MakeGatherSalesVolumesEndpoint(s service.GatherService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GatherSalesVolumesRequest)
		num, err := s.GatherSalesVolumes(ctx, req.TargetID, req.StartTime, req.EndTime)
		return GatherSalesVolumesResponse{
			Err: err,
			Num: num,
		}, nil
	}
}

// Failed implements Failer.
func (r GatherSalesVolumesResponse) Failed() error {
	return r.Err
}

// GatherCommentsNumberRequest collects the request parameters for the GatherCommentsNumber method.
type GatherCommentsNumberRequest struct {
	TargetID  int64  `json:"target_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// GatherCommentsNumberResponse collects the response parameters for the GatherCommentsNumber method.
type GatherCommentsNumberResponse struct {
	Num int64 `json:"num"`
	Err error `json:"err"`
}

// MakeGatherCommentsNumberEndpoint returns an endpoint that invokes GatherCommentsNumber on the service.
func MakeGatherCommentsNumberEndpoint(s service.GatherService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GatherCommentsNumberRequest)
		num, err := s.GatherCommentsNumber(ctx, req.TargetID, req.StartTime, req.EndTime)
		return GatherCommentsNumberResponse{
			Err: err,
			Num: num,
		}, nil
	}
}

// Failed implements Failer.
func (r GatherCommentsNumberResponse) Failed() error {
	return r.Err
}

// GatherViewsRequest collects the request parameters for the GatherViews method.
type GatherViewsRequest struct {
	TargetID  int64  `json:"target_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// GatherViewsResponse collects the response parameters for the GatherViews method.
type GatherViewsResponse struct {
	Num int64 `json:"num"`
	Err error `json:"err"`
}

// MakeGatherViewsEndpoint returns an endpoint that invokes GatherViews on the service.
func MakeGatherViewsEndpoint(s service.GatherService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GatherViewsRequest)
		num, err := s.GatherViews(ctx, req.TargetID, req.StartTime, req.EndTime)
		return GatherViewsResponse{
			Err: err,
			Num: num,
		}, nil
	}
}

// Failed implements Failer.
func (r GatherViewsResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// GatherSalesVolumes implements Service. Primarily useful in a client.
func (e Endpoints) GatherSalesVolumes(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	request := GatherSalesVolumesRequest{
		EndTime:   endTime,
		StartTime: startTime,
		TargetID:  targetID,
	}
	response, err := e.GatherSalesVolumesEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GatherSalesVolumesResponse).Num, response.(GatherSalesVolumesResponse).Err
}

// GatherCommentsNumber implements Service. Primarily useful in a client.
func (e Endpoints) GatherCommentsNumber(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	request := GatherCommentsNumberRequest{
		EndTime:   endTime,
		StartTime: startTime,
		TargetID:  targetID,
	}
	response, err := e.GatherCommentsNumberEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GatherCommentsNumberResponse).Num, response.(GatherCommentsNumberResponse).Err
}

// GatherViews implements Service. Primarily useful in a client.
func (e Endpoints) GatherViews(ctx context.Context, targetID int64, startTime string, endTime string) (num int64, err error) {
	request := GatherViewsRequest{
		EndTime:   endTime,
		StartTime: startTime,
		TargetID:  targetID,
	}
	response, err := e.GatherViewsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GatherViewsResponse).Num, response.(GatherViewsResponse).Err
}
