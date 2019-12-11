package main

import (
	"fmt"
	"log"
	"time"

	pb "practical_grpc/client_streaming/server/proto"
)

// DatabaseService is an implementation of the Database service in database.proto
type DatabaseService struct{}

// Search returns a stream of matching search results
func (db *DatabaseService) Search(r *pb.SearchRequest, s pb.Database_SearchServer) error {
	responses := []string{
		"Highest ranked content",
		"Some ranked content",
		"Some ranked content",
		"Lowest ranked content",
	}

	fmt.Printf("Server reached, after two seconds, results will be return in a stream: %v\n", time.Now().Unix())
	time.Sleep(2 * time.Second)

	for idx, resp := range responses {
		result := &pb.SearchResponse{MatchedTerm: r.Term, Rank: int32(idx + 1), Content: resp}

		fmt.Printf("Response delayed 1s: %v\n", time.Now().Unix())
		time.Sleep(1 * time.Second)

		if err := s.Send(result); err != nil {
			log.Printf("Error sending message to the client: %v", err)
			return err
		}
	}

	return nil
}
