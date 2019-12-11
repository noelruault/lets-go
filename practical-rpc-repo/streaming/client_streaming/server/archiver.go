package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"

	pb "practical_grpc/client_streaming/server/proto"
)

// ArchiverService is an implementation of the Archiver service in archiver.proto
type ArchiverService struct{}

// Zip generates a zip file from the streamed request messages
//
// The basic idea here is to start a loop that will end when either the client
// has sent all of its messages, resulting in an io.EOF error when you call
// Recv, or if any other error occurs while streaming in the messages.
func (as *ArchiverService) Zip(stream pb.Archiver_ZipServer) error {
	buf := new(bytes.Buffer)
	zf := zip.NewWriter(buf)

	for {
		// get or wait for the next request object in the stream:
		// "Each message you receive from the client is added to the writer."
		req, err := stream.Recv()

		// we're done, send the zip file:
		// "Once you have all of the messages from the client, send a response
		// back with the compressed bytes from the writer."
		if err == io.EOF {
			if err = zf.Close(); err != nil {
				log.Printf("Error creating the zip file: %v", err)
				return err
			}

			return stream.SendAndClose(&pb.ZipResponse{
				ZippedContents: buf.Bytes(),
			})
		}

		// an error occured when getting the request object
		if err != nil {
			log.Printf("Error reading the request: %v", err)
			return err
		}

		f, err := zf.Create(req.FileName)
		if err != nil {
			log.Printf("Error creating the zip file entry: %v", err)
			return err
		}

		if _, err := f.Write(req.Contents); err != nil {
			log.Printf("Error writing zip file: %v", err)
			return err
		}
	}
}
