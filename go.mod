module CPS831-Final

go 1.16

require (
	github.com/golang/protobuf v1.5.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.23.0
	google.golang.org/protobuf v1.26.0 // indirect
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
