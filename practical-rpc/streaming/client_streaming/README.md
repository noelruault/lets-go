# Client streaming

## Instructions

This code needs to be clone at `$GOPATH/src/practical_grpc/streaming` to work.

Run the following to create server/proto/archiver.pb.go:

```sh
cd $GOPATH/src/practical_grpc/client_streaming
docker run --rm -v $(pwd):/defs \
  namely/protoc-all:1.9 -d proto -l go -o server/proto
```

You’ll need a new directory for the Python code and need to include the grpcio package in order to make gRPC calls. To do that, run the following commands:

```sh
cd $GOPATH/src/practical_grpc/client_streaming
mkdir -p archiver/proto && cd archiver
pip install --user pipenv
pipenv install grpcio # pipenv install grpcio==1.11.0
cd ../
```

generate the Python client stubs:

```sh
cd $GOPATH/src/practical_grpc/client_streaming
docker run --rm -v $(PWD):/defs \
  namely/protoc-all:1.9 -f proto/archiver.proto -l python -o archiver
```

To configure the archiver:

```sh
cd $GOPATH/src/practical_grpc/client_streaming/archiver

pip uninstall pipenv
pip install --user pipenv

pipenv sync --sequential
```

## Run

To run the server:

```sh
cd $GOPATH/src/practical_grpc/client_streaming
go run server/main.go server/archiver.go server/database.go
```

To run the archiver:

`pipenv run python3 main.py <path/to/file>`

```sh
cd $GOPATH/src/practical_grpc/client_streaming/archiver
pipenv run python3 main.py files/compress_me.txt
```

## Extended docs

* [Github Project, streaming.](https://github.com/backstopmedia/gRPC-book-example/tree/master/chapters/streaming)
* [Practical gRPC book, chapter 6.](https://learning.oreilly.com/library/view/practical-grpc/9781939902580/ch06.html#idm140587010935408)

<!-- I dunnno if this worked when installing dependencies...
  pip install -e git+https://github.com/pypa/pipenv.git@master#egg=pipenv
  pipenv install grpcio==1.11.0
--->
