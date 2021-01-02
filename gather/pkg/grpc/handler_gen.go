// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package grpc

import (
	grpc "github.com/go-kit/kit/transport/grpc"
	endpoint "github.com/gw123/echo-app/gather/pkg/endpoint"
	pb "github.com/gw123/echo-app/gather/pkg/grpc/pb"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	gatherSalesVolumes   grpc.Handler
	gatherCommentsNumber grpc.Handler
	gatherViews          grpc.Handler
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.GatherServer {
	return &grpcServer{
		gatherCommentsNumber: makeGatherCommentsNumberHandler(endpoints, options["GatherCommentsNumber"]),
		gatherSalesVolumes:   makeGatherSalesVolumesHandler(endpoints, options["GatherSalesVolumes"]),
		gatherViews:          makeGatherViewsHandler(endpoints, options["GatherViews"]),
	}
}