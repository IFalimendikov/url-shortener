package grpchandler

import (
	"context"
	"errors"
	"url-shortener/internal/proto"
	"url-shortener/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Creates a shortened version of a URL provided in JSON format
func (g *GRPCHandler) ShortenURL(ctx context.Context, in *proto.ShortenURLRequest) (*proto.ShortenURLResponse, error) {
	var response *proto.ShortenURLResponse

	url := in.Url
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "URL is empty")
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

	shortURL, err := g.service.SaveURL(ctx, url, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorDuplicate) {
			return nil, status.Error(codes.AlreadyExists, "URL already exists")
		}
		return nil, status.Error(codes.Internal, "Couldn't encode URL")
	}

	resURL := g.cfg.BaseURL + "/" + string(shortURL)

	response = &proto.ShortenURLResponse{
		Result: resURL,
	}

	return response, nil
}
