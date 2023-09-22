package main

import (
	"net"
	"os"

	image "github.com/null-channel/eddington/container-builder/internal/containers/buildpack"
	"github.com/null-channel/eddington/container-builder/server"
	"github.com/null-channel/eddington/container-builder/utils"
	"github.com/null-channel/eddington/proto/container"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const DockerAPIVersion = "1.43"

func main() {

	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	grpc := grpc.NewServer()
	reflection.Register(grpc)

	l := zerolog.New(os.Stdout).With().Timestamp().Str("component", "server").Logger()

	db, err := utils.NewDB("")
	if err != nil {
		l.Err(err).Str("failed to create db", "")
	}

	registry := os.Getenv("REGISTRY_URL")
	if registry == "" {
		l.Error().Err(err).Msg("unable to create registry ")
	}

	builder, err := image.NewBuilder(db, registry)
	if err != nil {
		l.Error().Err(err).Msg("unable to create builder")
	}

	server, err := server.NewServer(db, builder, &l)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create server")
	}
	container.RegisterContainerServiceServer(grpc, server)

	log.Info().Msg("starting server on port 4040")

	err = grpc.Serve(lis)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to serve")
	}

}
