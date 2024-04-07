package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/goodbaikin/ebjs/api"
	"github.com/goodbaikin/ebjs/encode"
	"github.com/goodbaikin/ebjs/server"
	"google.golang.org/grpc"
)

const (
	port = 3000
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug")
	vcodec := flag.String("vcodec", "libx264", "video codec for encoding")
	acodec := flag.String("acodec", "libfdk_aac", "audio codec for encoding")
	flag.Parse()

	encodeOptions := encode.NewEncodeOptions()
	opts := []encode.EncodeOption{encode.WithVCodec(*vcodec)}
	opts = append(opts, encode.WithACodec(*acodec))
	for _, opt := range opts {
		opt(encodeOptions)
	}

	s := server.NewServer(flag.Arg(0), flag.Arg(1), encodeOptions, *debug)
	grpcServer := grpc.NewServer()
	api.RegisterEncoderServer(grpcServer, &s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	grpcServer.Serve(listener)
}
