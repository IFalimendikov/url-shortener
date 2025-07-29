package grpchandler

import (
	"context"
	"url-shortener/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Check if database connection is alive
func (g *GRPCHandler) PingDB(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	live := g.service.PingDB()
	if !live {
		return nil, status.Error(codes.Internal, "Database connection is not alive")
	}

	return &proto.PingResponse{
		Status: true,
	}, nil
}
