package stratego_rpc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type dummy struct{}

func validateUserId(userId string) error {
	if userId == "" {
		return fmt.Errorf("invalid user-id provided")
	}
	return nil
}

func getUserIdFromContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); !ok || md["x-stratego-user-id"] == nil || len(md["x-stratego-user-id"]) == 0 {
		return "", fmt.Errorf("no user-id provided")
	} else {
		userId := md["x-stratego-user-id"][0]
		if err := validateUserId(userId); err == nil {
			return userId, nil
		} else {
			return "", err
		}
	}
}

func requireUserIdDecorator[I interface{}, O interface{}](ctx context.Context, handler func(string, *I) (*O, error), input *I) (*O, error) {
	if userId, err := getUserIdFromContext(ctx); err == nil {
		return handler(userId, input)
	} else {
		log.Error().
			Err(err).
			Msg("failed to get user id")
		return nil, status.Errorf(codes.PermissionDenied, "invalid user id")
	}
}
