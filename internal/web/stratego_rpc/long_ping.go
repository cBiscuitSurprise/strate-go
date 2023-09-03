package stratego_rpc

import (
	"io"
	"time"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type longPingState struct {
	Prepend        string
	PreprendLittle bool
}

func (s *strateGoServer) LongPing(stream pb.StrateGo_LongPingServer) error {
	state := longPingState{
		Prepend:        "[p] ",
		PreprendLittle: true,
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if response, err := handleLongPingTermination(state); err != nil {
				return err
			} else {
				if err := stream.Send(response); err != nil {
					return err
				}
			}
			return nil
		} else if err != nil {
			return err
		}

		if response, err := handleLongPingRequest(in, state); err != nil {
			return err
		} else {
			if err := stream.Send(response); err != nil {
				return err
			}
		}

		if state.PreprendLittle {
			state = longPingState{
				Prepend:        "[o] ",
				PreprendLittle: false,
			}
		} else {
			state = longPingState{
				Prepend:        "[O] ",
				PreprendLittle: true,
			}
		}
	}
}

func handleLongPingRequest(request *pb.LongPingRequest, state longPingState) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   state.Prepend + request.Message,
		Games:     []*pb.Game{},
	}, nil
}

func handleLongPingTermination(state longPingState) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "[ng] byeeeeee 👋 ...",
		Games:     []*pb.Game{},
	}, nil
}
