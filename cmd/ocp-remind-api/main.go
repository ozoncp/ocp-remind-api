package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ocpremindapi "github.com/ozoncp/ocp-remind-api/internal/app/ocp-remind-api"
	"github.com/ozoncp/ocp-remind-api/pkg"
)

func LoadConfiguration(filePath string) error {
	openFile := func() error {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Println("error on file closing...")
			}
		}()
		return nil
	}
	for i := 0; i < 10; i++ {
		err := openFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func run() error {
	listen, err := net.Listen("tcp4", ":82")
	if err != nil {
		log.Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	pkg.RegisterRemindApiV1Server(s, ocpremindapi.NewRemindAPIV1())
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Err(err).Msg("failed to serve")
	}

	return nil
}

func main() {
	log.Printf("ocp remind api project")

	exitDone := &sync.WaitGroup{}
	exitDone.Add(1)

	if err := run(); err != nil {
		log.Err(err).Msg("failed to run")
	}
}
