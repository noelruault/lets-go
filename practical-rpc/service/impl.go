package service

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/noelruault/practical-rpc/proto"
)

// To start with, we'll hardcode the database of films.
var films = []*proto.Film{
	&proto.Film{
		Id:          "4",
		Title:       "A New Hope",
		Director:    "George Lucas",
		Producer:    "Gary Kurtz, Rick McCallum",
		ReleaseDate: toProto(1977, 5, 25),
	},
	&proto.Film{
		Id:          "5",
		Title:       "The Empire Strikes Back",
		Director:    "Irvin Kershner",
		Producer:    "Gary Kurtz, Rick McCallum",
		ReleaseDate: toProto(1980, 5, 17),
	},
	&proto.Film{
		Id:          "6",
		Title:       "Return of the Jedi",
		Director:    "Richard Marquand",
		Producer:    "Howard G. Kazanjian, George Lucas, Rick McCallum",
		ReleaseDate: toProto(1983, 5, 25),
	},
}

func toProto(year, month, day int) *timestamp.Timestamp {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

type StarfriendsImpl struct {
}

// GetFilm queries a film by ID or returns an error if not found.
func (s *StarfriendsImpl) GetFilm(ctx context.Context, req *proto.GetFilmRequest) (*proto.GetFilmResponse, error) {

	var film *proto.Film
	for _, f := range films {
		if f.Id == req.Id {
			film = f
			break
		}
	}
	if film == nil {
		return nil, status.Errorf(codes.NotFound, "no film with id %q", req.Id)
	}

	// GET HEADERS
	if reqHdrs, ok := metadata.FromIncomingContext(ctx); ok {
		// Must use all-lower-case keys to query metadata
		fmt.Printf("%+v", reqHdrs)
		if who, ok := reqHdrs["who"]; ok {
			// who is a slice of strings; just use the first
			log.Printf("Received request from %s", who[0])
		}
	}

	// SET HEADERS
	// For an unary RPC header are sent with every message and can be set in
	// the initial context, with grpc.SendHeader, grpc.SetHeader
	// and grpc.SetTrailer.
	respHdrs := metadata.New(map[string]string{
		"Who":     "starfriends-server",
		"Version": "v2",
	})
	grpc.SendHeader(ctx, respHdrs)

	start := time.Now()
	defer func() {
		respTrlrs := metadata.Pairs("duration", time.Since(start).String())
		grpc.SetTrailer(ctx, respTrlrs)
	}()

	return &proto.GetFilmResponse{Film: film}, nil
}

// ListFilms returns a list of all known films.
func (s *StarfriendsImpl) ListFilms(ctx context.Context,
	req *proto.ListFilmsRequest) (*proto.ListFilmsResponse, error) {

	// Note that for a stream there are only three times that headers can be
	// sent: in the context used to open the initial stream, via
	// grpc.SendHeader, and grpc.SetTrailer. It is not possible to set headers
	// on arbitrary messages in a stream.

	return &proto.ListFilmsResponse{Films: films}, nil
}

// compile-type check that our new type provides the
// correct server interface
var _ proto.StarfriendsServer = (*StarfriendsImpl)(nil)
