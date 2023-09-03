package web

import (
	"context"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RequestInfo struct {
	UserId    string
	RequestId string
}

// a UnaryServerInterceptor to log the function name
func UnaryRequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Warn().
			Str("FullMethod", info.FullMethod).
			Msgf("failed to get request metdata")
	}

	requestInfo := logServerCommand(md, info.FullMethod, "handling request")

	resp, err = handler(ctx, req)

	if err != nil {
		logError(err, info.FullMethod, requestInfo, "error handling request")
	}

	return resp, err
}

type streamRequestLoggerStreamWrapper struct {
	grpc.ServerStream

	Metadata   metadata.MD
	FullMethod string
}

func (s *streamRequestLoggerStreamWrapper) RecvMsg(m interface{}) error {
	info := logServerCommand(s.Metadata, s.FullMethod, "receiving request")
	if err := s.ServerStream.RecvMsg(m); err != nil {
		if err == io.EOF {
			logServerCommand(s.Metadata, s.FullMethod, "closed the stream (client)")
		} else {
			logError(err, s.FullMethod, info, "error sending response")
		}
		return err
	}
	return nil
}

func (s *streamRequestLoggerStreamWrapper) SendMsg(m interface{}) error {
	info := logServerCommand(s.Metadata, s.FullMethod, "sending response")
	if err := s.ServerStream.SendMsg(m); err != nil {
		if err == io.EOF {
			logServerCommand(s.Metadata, s.FullMethod, "closed the stream (server)")
		} else {
			logError(err, s.FullMethod, info, "error sending response")
		}
		return err
	}
	return nil
}

// a StreamServerInterceptor to log the function name
func StreamRequestLogger(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	fullMethod := info.FullMethod
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		log.Warn().
			Str("FullMethod", fullMethod).
			Msgf("failed to get request metdata")
	}

	wrapper := &streamRequestLoggerStreamWrapper{
		ServerStream: ss,

		Metadata:   md,
		FullMethod: fullMethod,
	}

	return handler(srv, wrapper)
}

func logServerCommand(md metadata.MD, fullMethod string, msg string, parts ...any) RequestInfo {
	output := RequestInfo{
		RequestId: strings.Join(md["x-request-id"], ", "),
		UserId:    strings.Join(md["x-stratego-user-id"], ", "),
	}

	log.Info().
		Str("FullMethod", fullMethod).
		Str("RequestId", output.RequestId).
		Str("UserId", output.UserId).
		Msgf(msg, parts...)

	return output
}

func logError(err error, fullMethod string, info RequestInfo, msg string) {
	log.Error().
		Err(err).
		Str("FullMethod", fullMethod).
		Str("RequestId", info.RequestId).
		Str("UserId", info.UserId).
		Msg(msg)
}
