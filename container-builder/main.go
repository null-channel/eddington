package main

import (
	"net"

	"github.com/null-channel/eddington/container-builder/server"
	"github.com/null-channel/eddington/proto/container"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	grpc := grpc.NewServer()
	reflection.Register(grpc)
	server, err := server.NewServer()
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
