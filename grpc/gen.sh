#!/bin/bash

protoc --go_out=plugins=grpc:. hb.proto
