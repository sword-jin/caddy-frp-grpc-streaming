package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/sword-jin/caddy-frp-grpc-streaming/proto"
	"github.com/sword-jin/caddy-frp-grpc-streaming/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle(protoconnect.NewServiceHandler(&handler{}))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	server.ListenAndServe()
}

type handler struct{}

var _ protoconnect.ServiceHandler = (*handler)(nil)

// Stream implements protoconnect.ServiceHandler.
func (h *handler) Stream(_ context.Context, _ *connect.Request[proto.Request], stream *connect.ServerStream[proto.Response]) error {
	for range 5 {
		if err := stream.Send(&proto.Response{
			Result: "Sreaming!",
		}); err != nil {
			return err
		}
	}
	return nil
}

// Unary implements protoconnect.ServiceHandler.
func (h *handler) Unary(context.Context, *connect.Request[proto.Request]) (*connect.Response[proto.Response], error) {
	return connect.NewResponse(&proto.Response{
		Result: "Unary!",
	}), nil
}
