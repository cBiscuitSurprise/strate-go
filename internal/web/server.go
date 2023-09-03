package web

import (
	"net"

	"github.com/cBiscuitSurprise/strate-go/internal/web/stratego_rpc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	opts = append(opts, grpc.StreamInterceptor(StreamRequestLogger))

	grpcServer := stratego_rpc.NewServer(opts)
	grpcServer.Serve(lis)
}
