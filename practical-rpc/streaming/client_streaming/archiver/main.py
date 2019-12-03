#!/usr/bin/env python3

import sys
import os
import time

import grpc

from proto import archiver_pb2
from proto import archiver_pb2_grpc as rpc


def _read(path):
    with open(path, 'rb') as stream:
        return archiver_pb2.ZipRequest(file_name=path, contents=stream.read())

# The archive function takes a set of file paths supplied in the command line
# and creates an iterator of ZipRequest messages to be sent to the server.
def archive(file_paths):
    """Send each file to the archiver service to be zipped up.
    The resulting zip file will be written to compressed.zip in
    the current working directory.
    """
    channel = grpc.insecure_channel('localhost:8080')
    stub = rpc.ArchiverStub(channel)

    files = iter([_read(path) for path in file_paths])
    future = stub.Zip.future(files)

    # calling future.result() will block until the server has sent its response
    # we'll just let any errors go to the console
    response = future.result()

    file_name = 'compressed.zip'
    with open(file_name, 'wb') as stream:
        # stream.write(response.zipped_contents.decode("utf-8"))
        stream.write(response.zipped_contents)

    print("Wrote " + file_name + " to the current directory.")

    print("The file will be removed in 10 seconds...")
    time.sleep(5)
    os.remove(file_name)
    print("File Removed!")


if __name__ == '__main__':
    archive(sys.argv[1:])
