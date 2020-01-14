package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/proto"
	"github.com/noelruault/programming-training/practical-rpc/golang-basic-gRPC-example/store"
	"golang.org/x/net/context"
)

func toGtime(ts *timestamp.Timestamp) time.Time {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		panic(err)
	}
	return t
}

func toProto(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

type StarFriendsImpl struct {
	store *store.InMemStore
}

func NewStarFriendsImpl(store *store.InMemStore) *StarFriendsImpl {
	return &StarFriendsImpl{
		store: store,
	}
}

// GetFilm queries a film by ID or returns an error if not found.
func (s *StarFriendsImpl) GetFilm(ctx context.Context, req *proto.GetFilmRequest) (
	*proto.GetFilmResponse, error) {

	f, err := s.store.Get(req.Id)
	if err != nil { // err = "not_found"
		if err.Error() == "not_found" {
			// This will be executed in the server side
			fmt.Println("The film wasn't found.")
		}
	}

	return &proto.GetFilmResponse{
		Film: &proto.Film{
			Id:          f.ID,
			Title:       f.Title,
			Director:    f.Director,
			Producer:    f.Producer,
			ReleaseDate: toProto(f.ReleaseDate),
		}}, nil
}

// ListFilms returns a list of all known films.
func (s *StarFriendsImpl) ListFilms(ctx context.Context,
	req *proto.ListFilmsRequest) (*proto.ListFilmsResponse, error) {

	// Note that for a stream there are only three times that headers can be
	// sent: in the context used to open the initial stream, via
	// grpc.SendHeader, and grpc.SetTrailer. It is not possible to set headers
	// on arbitrary messages in a stream.

	var films []*proto.Film
	fs := s.store.List()
	for _, f := range fs {
		films = append(films, &proto.Film{
			Id:          f.ID,
			Title:       f.Title,
			Director:    f.Director,
			Producer:    f.Producer,
			ReleaseDate: toProto(f.ReleaseDate),
		})
	}
	return &proto.ListFilmsResponse{Films: films}, nil
}

// DeleteFilm tries to delete a film by ID or returns an error if not found.
func (s *StarFriendsImpl) DeleteFilm(ctx context.Context,
	req *proto.DeleteFilmRequest) (*empty.Empty, error) {

	err := s.store.Delete(req.Id)
	if err != nil {
		return nil, errors.New("500")
	}

	return &empty.Empty{}, nil
}

func (s *StarFriendsImpl) CreateFilm(ctx context.Context,
	req *proto.CreateFilmRequest) (*proto.CreateFilmResponse, error) {

	f := s.store.Create(store.Film{
		Title:       req.Film.Title,
		Director:    req.Film.Director,
		Producer:    req.Film.Producer,
		ReleaseDate: toGtime(req.Film.ReleaseDate),
	})

	// TODO(nruault): Error handling

	return &proto.CreateFilmResponse{
		Film: &proto.Film{
			Id:          f.ID,
			Title:       f.Title,
			Director:    f.Director,
			Producer:    f.Producer,
			ReleaseDate: toProto(f.ReleaseDate),
		}}, nil
}

func (s *StarFriendsImpl) UpdateFilm(ctx context.Context,
	req *proto.UpdateFilmRequest) (*proto.UpdateFilmResponse, error) {

	f, err := s.store.Update(store.Film{
		ID:          req.Film.Id,
		Title:       req.Film.Title,
		Director:    req.Film.Director,
		Producer:    req.Film.Producer,
		ReleaseDate: toGtime(req.Film.ReleaseDate),
	})
	if err != nil {
		fmt.Errorf("%v", err)
	}

	// TODO(nruault): Error handling

	return &proto.UpdateFilmResponse{
		Film: &proto.Film{
			Id:          f.ID,
			Title:       f.Title,
			Director:    f.Director,
			Producer:    f.Producer,
			ReleaseDate: toProto(f.ReleaseDate),
		}}, nil
}

// compile-type check that our new type provides the correct server interface
var _ proto.StarfriendsServer = (*StarFriendsImpl)(nil)
