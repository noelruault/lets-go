package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/proto"
	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/service"
	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/store"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 8080, "The port on which gRPC server will listen")
	flag.Parse()

	// We're not providing TLS options, so server will use plaintext.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fail(err)
	}
	fmt.Printf("Listening tcp connections on %v\n", lis.Addr())

	svr := grpc.NewServer()
	store := store.NewInMemStore()
	api := service.NewStarFriendsImpl(store)
	proto.RegisterStarfriendsServer(svr, api)

	// trap SIGINT / SIGTERM to exit cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down...")
		svr.GracefulStop()
	}()

	// finally, run the server
	if err := svr.Serve(lis); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
