package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/lib/pq" // here
	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func toProto(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "postgres://postgres:pass@localhost/filmstore?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM films")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	stub := proto.NewStarfriendsClient(conn)
	// Now we can use the stub to make RPCs

	// SET HEADERS
	// OPTION 1 - Pairs
	// With the metadata.Pairs, you can specify the same key multiple times to
	// provide multiple values.
	//
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

	var respHdrs, respTrlrs metadata.MD

	// LIST
	req := &proto.ListFilmsRequest{}
	resp, err := stub.ListFilms(ctx, req,
		grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))

	// GET
	// req := &proto.GetFilmRequest{Id: "c849f921-5488-4750-b989-70ce002eb572"}
	// resp, err := stub.GetFilm(ctx, req,
	// 	grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))

	// CREATE
	// req := &proto.CreateFilmRequest{
	// 	Film: &proto.Film{
	// 		Title:    "The Farewell",
	// 		Director: "Gregory Molina",
	// 		Producer: "The goat records",
	// 		ReleaseDate: toProto(
	// 			time.Date(2007, time.Month(4), 2, 0, 0, 0, 0, time.Local)),
	// 	},
	// }
	// resp, err := stub.CreateFilm(ctx, req,
	// 	grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))

	// UPDATE
	// req := &proto.UpdateFilmRequest{
	// 	Film: &proto.Film{
	// 		Id:       "c849f921-5488-4750-b989-70ce002eb572",
	// 		Title:    "A New Hope (remastered)",
	// 		Director: "George Lucas",
	// 		Producer: "Gary Kurtz, Rick McCallum",
	// 		ReleaseDate: toProto(
	// 			time.Date(1977, time.Month(5), 25, 0, 0, 0, 0, time.Local)),
	// 	},
	// }
	// resp, err := stub.UpdateFilm(ctx, req,
	// 	grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))

	// DELETE
	// req := &proto.DeleteFilmRequest{Id: "c849f921-5488-4750-b989-70ce002eb572"}
	// resp, err := stub.DeleteFilm(ctx, req,
	// 	grpc.Header(&respHdrs), grpc.Trailer(&respTrlrs))

	// ERROR HANDLING
	if err != nil {
		fmt.Fprintf(os.Stderr, "RPC failed: %v\n", err)
	} else {
		fmt.Println(resp)
	}

	fmt.Printf("Server sent headers: %v\n", respHdrs)
	fmt.Printf("Server sent trailers: %v\n", respTrlrs)

}
