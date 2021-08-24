package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	_, cancel := context.WithCancel(context.Background())
	var grp errgroup.Group

	listen, err := net.Listen("tcp4", ":82")
	if err != nil {
		log.Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	api, err := ocpremindapi.NewRemindAPIV1()
	if err != nil {
		cancel()
		return err
	}
	pkg.RegisterRemindApiV1Server(s, api)
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Err(err).Msg("failed to serve")
	}

	osSignal := <-c
	log.Info().Msgf("system syscall:%+v", osSignal)

	s.GracefulStop()

	cancel()

	if err = grp.Wait(); err != http.ErrServerClosed {
		log.Fatal().Msgf("server shutdown failed: %v", err)
	}

	return nil
}

func main() {
	log.Printf("ocp remind api project")
	//dbConn := db.Connect("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	//exitDone := &sync.WaitGroup{}
	//exitDone.Add(1)
	if err := run(); err != nil {
		log.Err(err).Msg("failed to run")
	}
}
