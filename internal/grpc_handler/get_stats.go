package grpc_handler

import (
	"context"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"url-shortener/internal/proto"
	"google.golang.org/grpc/metadata"
)

// Shows service stats
func (g *GRPCHandler) GetStats(ctx context.Context, in *proto.GetStatsRequest) (*proto.GetStatsResponse, error) {
	var response *proto.GetStatsResponse

	if g.cfg.TrustedSubnet == "" {
		return nil, status.Error(codes.Internal, "Trusted subnet is not configured")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata not found")
	}
	values := md.Get("X-Real-IP")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "X-Real-IP header not provided")
	}
	userIP := net.ParseIP(values[0])

	_, network, err := net.ParseCIDR(g.cfg.TrustedSubnet)
	if err != nil {
		return nil, status.Error(codes.Internal, "Can't parse CIDR")
	}

	access := network.Contains(userIP)
	if !access {
		return nil, status.Error(codes.PermissionDenied, "IP address is not trusted")
	}

	stats, err := g.service.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Stats not found")
	}

	response = &proto.GetStatsResponse{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}

	return response, nil
}
