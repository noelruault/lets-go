# Server streaming

## Instructions

This code needs to be clone at `$GOPATH/src/practical_grpc/server_streaming` to work.

You'll need Golang, Python, and Ruby.

* Run `make setup` to get all the necessary tools.

* Consider the possibility of installing grpc gem locally:

```sh
gem install grpc
```

Run the following to create server/proto.

```sh
cd $GOPATH/src/practical_grpc/server_streaming
docker run --rm -v $(pwd):/defs \
  namely/protoc-all:1.9 -d proto -l go -o server/proto
```

Run the following to create server/client.

```sh
cd $GOPATH/src/practical_grpc/server_streaming
docker run --rm -v $(pwd):/defs \
  namely/protoc-all:1.9 -f proto/database.proto -l ruby -o database
```

## Run

To run the server.

```sh
cd $GOPATH/src/practical_grpc/server_streaming
go run server/main.go server/database.go
```

To run the client.

```sh
cd $GOPATH/src/practical_grpc/server_streaming
ruby database/client.rb
```

## Extended docs

* [Github Project, streaming.](https://github.com/backstopmedia/gRPC-book-example/tree/master/chapters/streaming)
* [Practical gRPC book, chapter 6.](https://learning.oreilly.com/library/view/practical-grpc/9781939902580/ch06.html#idm140587011558096)
