package grpchandler

import (
	"context"
	"errors"
	"net/url"
	"url-shortener/internal/proto"
	"url-shortener/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Creates a shortened version of a provided URL
func (g *GRPCHandler) PostURL(ctx context.Context, in *proto.ShortenURLRequest) (*proto.ShortenURLResponse, error) {
	var response *proto.ShortenURLResponse

	urlStr := in.Url
	if urlStr == "" {
		return nil, status.Error(codes.InvalidArgument, "URL is empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, status.Error(codes.InvalidArgument, "Malformed URI")
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

	shortURL, err := g.service.SaveURL(ctx, urlStr, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorDuplicate) {
			resURL := g.cfg.BaseURL + "/" + string(shortURL)
			return &proto.ShortenURLResponse{
				Result: resURL,
			}, nil
		}
		return nil, status.Error(codes.Internal, "Couldn't encode URL")
	}

	resURL := g.cfg.BaseURL + "/" + string(shortURL)
	response = &proto.ShortenURLResponse{
		Result: resURL,
	}

	return response, nil
}
