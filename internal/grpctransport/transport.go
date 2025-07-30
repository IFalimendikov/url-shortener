package grpctransport

import (
	"context"
	"log/slog"

	"url-shortener/internal/grpchandler"
	pb "url-shortener/internal/proto"

	"google.golang.org/grpc"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is a type for keys used in context to avoid collisions
type contextKey string

const (
	userIDKey contextKey = "user_id"
)

// Transport handles gRPC transport layer operations including middleware and routing
type GRPCTransport struct {
	handler *grpchandler.GRPCHandler
	log     *slog.Logger
}

// New creates a new Transport instance with the provided configuration and handlers
func New(h *grpchandler.GRPCHandler, log *slog.Logger) *GRPCTransport {
	return &GRPCTransport{
		handler: h,
		log:     log,
	}
}

// Claims represents JWT claims structure
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// NewRouter creates and returns a new configured gRPC server with interceptors
func NewRouter(g *GRPCTransport) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(g.log),
			AuthInterceptor(),
		),
	)

	pb.RegisterURLShortenerServer(server, g.handler)

	return server
}

// LoggingInterceptor adds request logging that records method, duration, and status
func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		method := info.FullMethod

		resp, err := handler(ctx, req)

		latency := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Internal
			}
		}

		log.Info("request completed",
			"method", method,
			"duration", latency.String(),
			"status", code.String(),
		)

		return resp, err
	}
}

// AuthInterceptor handles JWT authentication via metadata
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to get metadata from context")
		}

		var userID string

		// Check for Authorization header (standard approach)
		authHeaders := md.Get("authorization")
		if len(authHeaders) > 0 {
			authHeader := authHeaders[0]
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString := authHeader[7:]

				claims := &Claims{}
				token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, status.Error(codes.InvalidArgument, "unexpected signing method")
					}
					return []byte("123"), nil
				})

				if err == nil && token.Valid {
					userID = claims.UserID
					ctx = context.WithValue(ctx, userIDKey, userID)
					return handler(ctx, req)
				}
			}
		}

		return nil, status.Error(codes.Unauthenticated, "valid authentication required")
	}
}
