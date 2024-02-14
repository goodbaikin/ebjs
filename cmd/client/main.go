package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/goodbaikin/ebjs/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	server := flag.String("server", "", "address of server")
	input := flag.String("input", "", "input file name")
	output := flag.String("output", "", "output file name")
	channelID := flag.Uint64("channelID", 0, "channel ID")
	isDualMonoMode := flag.Bool("isDualMonoMode", false, "whether the file is dual mono mode")
	flag.Parse()

	if *server == "" || *input == "" || *output == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *channelID == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conn, err := grpc.Dial(*server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		printErrorAndExit(err)
	}

	client := api.NewEncoderClient(conn)
	stream, err := client.Encode(context.Background(), &api.EncodeRequest{
		Input:          *input,
		Output:         *output,
		ChannelId:      *channelID,
		IsDualMonoMode: *isDualMonoMode,
	})
	if err != nil {
		printErrorAndExit(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			printErrorAndExit(err)
		}

		fmt.Println(resp.Progress)
	}
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
