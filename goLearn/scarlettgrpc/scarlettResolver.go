package scarlettgrpc

import "google.golang.org/grpc/resolver"

const (
	myScheme   = "scarlett"
	myEndpoint = "resolver.scarlett.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973"}

type scarlettResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}


