package main

import (
	"net"
	"os"

	"github.com/null-channel/eddington/container-builder/proto/container"
	"github.com/null-channel/eddington/container-builder/server"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

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

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugaredLogger := logger.Sugar()

	kubeConfig := os.Getenv("KUBECONFIG")

	clusterConfig, err := getClusterConfig(kubeConfig)

	if err != nil {
		sugaredLogger.Fatal(err)
		return
	}

	registry := os.Getenv("REGISTRY_URL")
	if registry == "" {
		sugaredLogger.Errorw("registry not configured", "error:", err)
	}

	server, err := server.NewServer(clusterConfig, sugaredLogger)
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

// getClusterConfig return the config for k8s
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
