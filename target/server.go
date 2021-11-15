package target

import (
	"context"
	"fmt"
)

var _ GreeterServer = &Server{}

type Server struct {
	UnimplementedGreeterServer
}

func (Server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: fmt.Sprintf("Hello there %s", req.Name),
	}, nil
}

