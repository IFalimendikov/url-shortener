package grpchandler

import (
	"context"
	"errors"
	"url-shortener/internal/proto"
	"url-shortener/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Retrieves and redirects to the original URL from a shortened URL ID
func (g *GRPCHandler) GetURL(ctx context.Context, in *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	var response *proto.GetURLResponse
	id := in.ShortUrl

	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "URL is empty")
	}

	url, err := g.service.GetURL(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorURLDeleted) {
			return nil, status.Error(codes.NotFound, "URL was deleted")
		}
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	response = &proto.GetURLResponse{
		OriginalUrl: url,
	}

	return response, nil
}
