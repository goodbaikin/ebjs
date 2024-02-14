package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/goodbaikin/ebjs/api"
	"github.com/goodbaikin/ebjs/server"
	"google.golang.org/grpc"
)

const (
	port = 3000
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug")
	flag.Parse()

	s := server.NewServer(flag.Arg(0), flag.Arg(1), *debug)
	grpcServer := grpc.NewServer()
	api.RegisterEncoderServer(grpcServer, &s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	grpcServer.Serve(listener)
}
