package web

import (
	"context"
	"net"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type ServerTlsOptions struct {
	CertFile string
	KeyFile  string
}

type ServerOptions struct {
	Origin     string
	TlsOptions *ServerTlsOptions
}

func Serve(params ServerOptions) {
	log.Info().Msgf("serving on '%s'", params.Origin)

	lis, err := net.Listen("tcp", params.Origin)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to open socket")
	}

	var opts []grpc.ServerOption

	if params.TlsOptions != nil {
		creds, err := credentials.NewServerTLSFromFile(
			params.TlsOptions.CertFile,
			params.TlsOptions.KeyFile,
		)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("failed to generate credentials")
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	opts = append(opts, grpc.UnaryInterceptor(UnaryRequestLogger))

	grpcServer := newServer(opts)
	grpcServer.Serve(lis)
}

// a UnaryServerInterceptor to log the function name
func UnaryRequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Warn().
			Str("FullMethod", info.FullMethod).
			Msgf("failed to get request metdata")
	}

	requestId := strings.Join(md["x-request-id"], ", ")

	log.Info().
		Str("FullMethod", info.FullMethod).
		Str("RequestId", requestId).
		Msgf("serving request")

	resp, err = handler(ctx, req)

	if err != nil {
		log.Error().
			Err(err).
			Str("FullMethod", info.FullMethod).
			Str("RequestId", requestId).
			Msg("error from handler")
	}

	return resp, err
}
