package grpc

import (
	"context"
	"errors"
	grpc "github.com/go-kit/kit/transport/grpc"
	endpoint "github.com/gw123/echo-app/gather/pkg/endpoint"
	pb "github.com/gw123/echo-app/gather/pkg/grpc/pb"
	context1 "golang.org/x/net/context"
)

// makeGatherSalesVolumesHandler creates the handler logic
func makeGatherSalesVolumesHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GatherSalesVolumesEndpoint, decodeGatherSalesVolumesRequest, encodeGatherSalesVolumesResponse, options...)
}

// decodeGatherSalesVolumesResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GatherSalesVolumes request.
// TODO implement the decoder
func decodeGatherSalesVolumesRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Decoder is not impelemented")
}

// encodeGatherSalesVolumesResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeGatherSalesVolumesResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Encoder is not impelemented")
}
func (g *grpcServer) GatherSalesVolumes(ctx context1.Context, req *pb.GatherSalesVolumesRequest) (*pb.GatherSalesVolumesReply, error) {
	_, rep, err := g.gatherSalesVolumes.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GatherSalesVolumesReply), nil
}

// makeGatherCommentsNumberHandler creates the handler logic
func makeGatherCommentsNumberHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GatherCommentsNumberEndpoint, decodeGatherCommentsNumberRequest, encodeGatherCommentsNumberResponse, options...)
}

// decodeGatherCommentsNumberResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GatherCommentsNumber request.
// TODO implement the decoder
func decodeGatherCommentsNumberRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Decoder is not impelemented")
}

// encodeGatherCommentsNumberResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeGatherCommentsNumberResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Encoder is not impelemented")
}
func (g *grpcServer) GatherCommentsNumber(ctx context1.Context, req *pb.GatherCommentsNumberRequest) (*pb.GatherCommentsNumberReply, error) {
	_, rep, err := g.gatherCommentsNumber.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GatherCommentsNumberReply), nil
}

// makeGatherViewsHandler creates the handler logic
func makeGatherViewsHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GatherViewsEndpoint, decodeGatherViewsRequest, encodeGatherViewsResponse, options...)
}

// decodeGatherViewsResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GatherViews request.
// TODO implement the decoder
func decodeGatherViewsRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Decoder is not impelemented")
}

// encodeGatherViewsResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeGatherViewsResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Gather' Encoder is not impelemented")
}
func (g *grpcServer) GatherViews(ctx context1.Context, req *pb.GatherViewsRequest) (*pb.GatherViewsReply, error) {
	_, rep, err := g.gatherViews.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GatherViewsReply), nil
}
