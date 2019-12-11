package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	pb "practical_grpc/bidi_streaming/server/proto"
)

type TokenizerService struct{}

func (ts *TokenizerService) Tokenize(s pb.Tokenizer_TokenizeServer) error {

	// for loop that breaks when the client has finished sending messages,
	// or when an error occurs.
	for {
		req, err := s.Recv()
		if err == io.EOF {
			return nil
		}

		// an error occured when getting the request object
		if err != nil {
			log.Printf("Got an error receiving from the client: %v", err)
			return err
		}

		rdr := bufio.NewReader(bytes.NewReader(req.FileContents))
		scanner := bufio.NewScanner(rdr)
		scanner.Split(bufio.ScanWords)

		// Each message that’s received the service creates a map of words and
		// their counts. This map is then sent back to the client for
		// aggregation/reduction purposes.
		results := &pb.TokenizeResponse{
			Words: make(map[string]int64),
		}

		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())

			fmt.Printf("Response delayed 0.1s: %v\n", time.Now().Unix())
			time.Sleep(100 * time.Millisecond)

			if _, ok := results.Words[word]; ok {
				results.Words[word]++
			} else {
				results.Words[word] = 1
			}

		}

		if err = s.Send(results); err != nil {
			log.Printf("Got an error sending to the client: %v", err)
			return err
		}

		fmt.Printf("END OF STREAM\n")
	}
}
