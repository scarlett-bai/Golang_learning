protoc13 -I=. --go_out=plugins=grpc,paths=source_relative:gen/go trip.proto
protoc13 -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go trip.proto