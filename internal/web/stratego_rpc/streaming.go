package stratego_rpc

import (
	"io"
)

type RequestHandler[I interface{}, O interface{}, S interface{}] func(request *I, state *S) (response *O, err error)

type GrpcServer[I interface{}, O interface{}] interface {
	Send(*O) error
	Recv() (*I, error)
}

type StreamingRequestHandler[I interface{}, O interface{}, S interface{}] struct {
	stream    GrpcServer[I, O]
	state     *S
	process   RequestHandler[I, O, S]
	terminate RequestHandler[I, O, S]
}

func (h *StreamingRequestHandler[I, O, S]) Listen() error {
	for {
		in, err := h.stream.Recv()
		if err == io.EOF {
			return h.handleRequest(in, h.terminate)
		} else if err != nil {
			return err
		}

		if err := h.handleRequest(in, h.process); err != nil {
			return err
		}
	}
}

func (h *StreamingRequestHandler[I, O, S]) handleRequest(request *I, handler RequestHandler[I, O, S]) error {
	if response, err := handler(request, h.state); err != nil {
		return err
	} else {
		if err := h.stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}
