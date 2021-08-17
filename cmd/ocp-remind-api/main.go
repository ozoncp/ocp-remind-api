package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc" //nolint:gci

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

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
)

func runGrpc() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pkg.RegisterRemindApiV1Server(s, ocpremindapi.NewRemindAPIV1())

	if err := s.Serve(listen); err != nil {
		log.Fatal("failed to serve: %v", err)
	}

	return nil
}

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pkg.RegisterRemindApiV1Server(s, ocpremindapi.NewRemindAPIV1())

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func runJSON() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pkg.RegisterRemindApiV1HandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("ocp remind api project")

	go runJSON()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
