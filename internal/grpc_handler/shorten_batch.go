package grpc_handler

import (
	"context"
	"url-shortener/internal/models"
	"url-shortener/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
)

// ShortenBatch handles batch URL shortening requests via gRPC
func (g *GRPCHandler) ShortenBatch(ctx context.Context, in *proto.BatchRequest) (*proto.BatchResponse, error) {
	var response *proto.BatchResponse

	urls := in.Urls
	if len(urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty request")
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

	var req []models.BatchUnitURLRequest
	for _, item := range urls {
		req = append(req, models.BatchUnitURLRequest{
			ID:  item.CorrelationId,
			URL: item.OriginalUrl,
		})
	}

	var res []models.BatchUnitURLResponse
	err := g.service.ShortenBatch(ctx, userID, req, &res)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error processing batch")
	}

	response = &proto.BatchResponse{
		Urls: make([]*proto.BatchUnitURLResponse, 0, len(res)),
	}

	for _, item := range res {
		response.Urls = append(response.Urls, &proto.BatchUnitURLResponse{
			CorrelationId: item.ID,
			ShortUrl:      g.cfg.BaseURL + "/" + string(item.Short),
		})
	}

	return response, nil
}
