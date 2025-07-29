package grpc_handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"url-shortener/internal/proto"
)

// Delete multiple URLs for a specific user
func (g *GRPCHandler) DeleteURLs(ctx context.Context, in *proto.DeleteURLsRequest) (*proto.DeleteURLsResponse, error) {
	urls := in.Urls

	if len(urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty or malformed body sent!")
	}

	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("user_id")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "User ID not provided")
		}
		userID = values[0]
	} else {
		return nil, status.Error(codes.Unauthenticated, "Metadata not found")
	}

	err := g.service.DeleteURLs(urls, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error deleting URLs")
	}
	
	return &proto.DeleteURLsResponse{
		Status: true,
	}, nil
}
