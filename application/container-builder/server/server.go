package server

import (
	"context"
	"strings"

	image "github.com/null-channel/eddington/application/container-builder/internal/containers/buildpack"
	"github.com/null-channel/eddington/proto/container"
)

type Server struct {
	container.UnimplementedContainerServiceServer
	builder *image.Builder
}

func (s *Server) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.CreateContainerReply, error) {
	buildinfo, err := s.builder.GetBuildPackInfo(req.Type)
	if err != nil {
		return nil, err
	}

	opts := image.BuildOpt{
		RepoURL: req.RepoURL,
		// TODO: ImageName:
		ImageName: strings.TrimPrefix(req.RepoURL, "https://github.com/"),
		BuildPack: buildinfo.BuildPack,
		Builder:   buildinfo.Builder,
	}

	err = s.builder.CreateImage(opts)
	if err != nil {
		return nil, err
	}
	return &container.CreateContainerReply{
		Message: "success",
	}, nil
}

func (s *Server) ImageStatus(ctx context.Context, req *container.ContainerID) (*container.ContainerStatus, error) {

	return nil, nil
}
