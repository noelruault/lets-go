package main

import (
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/noelruault/practical-rpc/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	stub := proto.NewStarfriendsClient(conn)
	// Now we can use the stub to make RPCs

	// SET HEADERS
	// OPTION 1 - Pairs
	// ctx := metadata.NewOutgoingContext(
	// 	context.Background(),
	// 	metadata.Pairs(
	// 		"Who", "starfiends-go-client",
	// 		"version", "v1",
	// 	),
	// )

	// OPTION 2 - New
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.New(map[string]string{
			"Who":     "starfiends-go-client",
			"version": "v1",
		}),
	)

	/*
		req := &proto.GetFilmRequest{Id: "4"}
		resp, err := stub.GetFilm(ctx, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RPC failed: %v\n", err)
		} else {
			fmt.Println(resp)
		}
	*/

	// respHdrs := metadata.New(map[string]string{
	// 	"Who":     "starfriends-server",
	// 	"Version": "v1",
	// })
	// grpc.SetHeader(ctx, respHdrs)

	// start := time.Now()
	// defer func() {
	// 	respTrlrs := metadata.Pairs("duration", time.Since(start).String())
	// 	grpc.SetTrailer(ctx, respTrlrs)
	// }()

	// We'll make another request and also print the response metadata
	req := &proto.GetFilmRequest{Id: "4"}
	var respHdrs, respTrlrs metadata.MD
	// var respTrlrs metadata.MD
	resp, err := stub.GetFilm(ctx, req,
		grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "RPC failed: %v\n", err)
	} else {
		fmt.Println(resp)
	}

	fmt.Printf("Server sent headers: %v\n", respHdrs)
	fmt.Printf("Server sent trailers: %v\n", respTrlrs)
}
