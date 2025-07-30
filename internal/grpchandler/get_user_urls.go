package grpchandler

import (
	"url-shortener/internal/models"
	"url-shortener/internal/proto"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Retrieves all URLs associated with the authenticated user
func (g *GRPCHandler) GetUserURLs(ctx context.Context, in *proto.GetUserURLsRequest) (*proto.GetUserURLsResponse, error) {
	response := &proto.GetUserURLsResponse{
		Urls: make([]*proto.UserURLResponse, 0),
	}
	var urls []models.UserURLResponse

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

	err := g.service.GetUserURLs(ctx, userID, &urls)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error retrieving user URLs")
	}

	if len(urls) == 0 {
		return nil, status.Error(codes.NotFound, "No URLs found for the user")
	}

	for i, url := range urls {
		response.Urls[i] = &proto.UserURLResponse{
			ShortUrl:    g.cfg.BaseURL + "/" + url.ShortURL,
			OriginalUrl: url.OriginalURL,
		}
	}

	return response, nil
}
