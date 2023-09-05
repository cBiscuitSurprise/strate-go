package stratego_rpc

import (
	"time"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type longPingState struct {
	Prepend        string
	PreprendLittle bool
}

func (s *strateGoServer) LongPing(stream pb.StrateGo_LongPingServer) error {
	state := &longPingState{
		Prepend:        "[p] ",
		PreprendLittle: true,
	}

	handler := StreamingRequestHandler[pb.LongPingRequest, pb.Pong, longPingState]{
		stream:    stream,
		state:     state,
		process:   handleLongPingRequest,
		terminate: handleLongPingTermination,
	}

	return handler.Listen()
}

func handleLongPingRequest(request *pb.LongPingRequest, state *longPingState) (*pb.Pong, error) {
	output := &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   state.Prepend + request.Message,
		Games:     []*pb.Game{},
	}

	if state.PreprendLittle {
		state.Prepend = "[o] "
		state.PreprendLittle = false
	} else {
		state.Prepend = "[O] "
		state.PreprendLittle = true
	}

	return output, nil
}

func handleLongPingTermination(request *pb.LongPingRequest, state *longPingState) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "[ng] byeeeeee 👋 ...",
		Games:     []*pb.Game{},
	}, nil
}
