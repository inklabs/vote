#!/bin/bash

set -ex

sudo apt-get update
sudo apt-get install -y protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

python3 -m pip install --upgrade pip
pip install grpcio-tools
