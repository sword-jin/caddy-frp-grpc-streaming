package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/sword-jin/caddy-frp-grpc-streaming/proto"
	"github.com/sword-jin/caddy-frp-grpc-streaming/proto/protoconnect"
)

func main() {
	httpClient := &http.Client{}
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.TLSClientConfig.InsecureSkipVerify = true
	httpClient.Transport = httpTransport

	ctx := context.Background()
	client := protoconnect.NewServiceClient(httpClient, "https://grpc-server.example.local", connect.WithGRPC())

	println("---------- TEST UNARY ----------")
	resp, err := client.Unary(ctx, connect.NewRequest(&proto.Request{Id: 1}))
	if err != nil {
		log.Fatalf("failed to call Unary: %v", err)
	}
	println("---------- UNARY RESPONSE ----------")
	println(resp.Msg.Result)

	println("---------- TEST STREAM ----------")
	stream, err := client.Stream(ctx, connect.NewRequest(&proto.Request{Id: 1}))
	if err != nil {
		log.Fatalf("failed to call Stream: %v", err)
	}
	defer stream.Close()

	for stream.Receive() {
		resp := stream.Msg()
		if err != nil {
			log.Fatalf("failed to receive stream response: %v", err)
		}
		fmt.Printf("stream response: %v\n", resp.Result)
	}

	if err := stream.Err(); err != nil {
		log.Fatalf("stream error: ", err)
	}
	println("---------- STREAM RESPONSE ----------")
}
