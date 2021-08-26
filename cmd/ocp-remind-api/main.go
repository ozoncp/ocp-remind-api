package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ocpremindapi "github.com/ozoncp/ocp-remind-api/internal/app/ocp-remind-api"
	"github.com/ozoncp/ocp-remind-api/internal/metrics"
	"github.com/ozoncp/ocp-remind-api/internal/tracer"
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
	var _ pkg.RemindApiV1Server = (*ocpremindapi.RemindAPIV1)(nil)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("jaeger", promhttp.Handler())

	metricsSrv := &http.Server{
		Addr:    "6831",
		Handler: mux,
	}

	metrics.CreateMetrics()

	tr := tracer.InitTracer("reminds")

	listen, err := net.Listen("tcp4", ":82") //nolint:gosec
	if err != nil {
		log.Err(err).Msg("failed to listen")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("REMINDS_DB_URL"))
	if err != nil {
		log.Err(err).Msg("Unable to connect to db")
		return err
	}
	defer func(ctx context.Context, conn *pgx.Conn) {
		err := conn.Close(ctx)
		if err != nil {
			log.Err(err).Msg("Unable to close connection to db")
		}
	}(context.Background(), conn)

	s := grpc.NewServer()
	api, err := ocpremindapi.NewRemindAPIV1(conn)
	if err != nil {
		return err
	}
	pkg.RegisterRemindApiV1Server(s, api)
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Err(err).Msg("failed to serve")
	}

	osSignal := <-c
	log.Info().Msgf("system syscall:%+v", osSignal)
	err = metricsSrv.Close()
	if err != nil {
		log.Err(err).Msg("Error on close metrics server")
	}
	s.GracefulStop()
	err = tr.Close()
	if err != nil {
		log.Err(err).Msg("Error on close tracer server")
	}

	return nil
}

func main() {
	log.Printf("ocp remind api project")
	if err := run(); err != nil {
		log.Err(err).Msg("failed to run")
	}
}
