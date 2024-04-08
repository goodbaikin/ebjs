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
	hwaccel := flag.String("hwaccel", "", "use hardware acceleration")
	hwaccelOutputFormat := flag.String("hwaccel_output_format", "", "output format for hwaccel")
	vaapiDevice := flag.String("vaapi_device", "", "set vaapi device")
	vf := flag.String("vf", "", "set video filter")
	flag.Parse()

	encodeOptions := encode.NewEncodeOptions()
	opts := []encode.EncodeOption{
		encode.WithVCodec(*vcodec),
		encode.WithACodec(*acodec),
		encode.WithVF(*vf),
	}
	if *hwaccel != "" {
		opts = append(opts, encode.WithHWAccel(*hwaccel))
	}
	if *hwaccelOutputFormat != "" {
		opts = append(opts, encode.WithHWAccelOutputFormat(*hwaccelOutputFormat))
	}
	if *vaapiDevice != "" {
		opts = append(opts, encode.WithVAAPIDevice(*vaapiDevice))
	}
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
