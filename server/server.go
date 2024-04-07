package server

import (
	"fmt"
	"os"

	"github.com/goodbaikin/ebjs/api"
	"github.com/goodbaikin/ebjs/encode"
)

type Server struct {
	encoder encode.Encoder
	debug   bool
	api.UnimplementedEncoderServer
}

func NewServer(basedir, recordedDir string, encodeOption *encode.EncodeOptions, debug bool) Server {
	s := Server{
		encoder: encode.NewEncoder(basedir, recordedDir, encodeOption),
		debug:   debug,
	}
	return s
}

func (s *Server) Encode(r *api.EncodeRequest, stream api.Encoder_EncodeServer) error {
	logger := func(msg string) {
		if s.debug {
			fmt.Println(msg)
		}
		stream.Send(&api.EncodeProgress{
			Progress: msg,
		})
	}

	if err := s.encoder.Encode(r.Input, r.Output, r.ChannelId, r.IsDualMonoMode, logger); err != nil {
		logger(err.Error())
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		return err
	}

	return nil
}
