# Client streaming

## Instructions

This code needs to be clone at `$GOPATH/src/practical_grpc/streaming` to work.

We’ll need to generate some server code with protoc:

```sh
cd $GOPATH/src/practical_grpc/bidi_streaming
docker run --rm -v $(pwd):/defs \
  namely/protoc-all:1.9 -d proto -l go -o server/proto
```

To generate the client stubs:

```sh
cd $GOPATH/src/practical_grpc/bidi_streaming
docker run --rm -v $(PWD):/defs \
  namely/protoc-all:1.9 -f proto/tokenizer.proto -l node -o tokenizer
```

## Run

To run the server:

```sh
cd $GOPATH/src/practical_grpc/bidi_streaming
go run server/main.go server/archiver.go server/tokenizer.go
```

To run the tokenizer:

`node index.js <path/to/file>`

```sh
cd $GOPATH/src/practical_grpc/bidi_streaming/tokenizer
node index.js package.json
```

## Extended docs

* [Github Project, streaming.](https://github.com/backstopmedia/gRPC-book-example/tree/master/chapters/streaming)
* [Practical gRPC book, chapter 6.](https://learning.oreilly.com/library/view/practical-grpc/9781939902580/ch06.html)
