package server

import (
	"context"
	"database/sql"
	"os"

	"strings"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/google/uuid"
	image "github.com/null-channel/eddington/application/container-builder/internal/containers/buildpack"
	"github.com/null-channel/eddington/application/container-builder/models"
	"github.com/null-channel/eddington/proto/container"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Server struct {
	container.UnimplementedContainerServiceServer
	builder *image.Builder
	log     zerolog.Logger
}

func NewServer() (*Server, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create db")
	}
	_, err = db.NewCreateTable().Model((*models.Build)(nil)).Exec(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "unable to migrate table")
	}
	packClient, err := pack.NewClient()
	if err != nil {
		return nil, err
	}
	builder, err := image.NewBuilder(db, packClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create builder")
	}
	l := zerolog.New(os.Stdout).With().Timestamp().Str("component", "server").Logger()
	return &Server{
		builder: builder,
		log:     l,
	}, nil

}

// CreateContainer maps to the CreateContainer RPC
func (s *Server) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.CreateContainerReply, error) {
	buildID := uuid.New().String()
	// create a build request in the db
	repo := strings.TrimPrefix(req.RepoURL, "https://github.com")
	err := s.builder.NewBuild(buildID, repo, req.CustomerID)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to create build request")
		return nil, errors.Wrap(err, "unable to create build , please try again")
	}

	// start build in a goroutine with the buildID and request
	go func(buildID string, req *container.CreateContainerRequest) {
		buildInfo, err := s.builder.GetBuildPackInfo(req.Type)
		if err != nil {
			s.log.Error().Err(err).Msg("unable to get buildpack info")
			return
		}

		opts := image.BuildOpt{
			BuildID:   buildID,
			RepoURL:   req.RepoURL,
			ImageName: strings.TrimPrefix(req.RepoURL, "https://github.com"),
			BuildPack: buildInfo.BuildPack,
			Builder:   buildInfo.Builder,
		}

		err = s.builder.CreateImage(opts)
		if err != nil {
			s.log.Error().Err(err).Msg("Build failed")
			return
		}

	}(buildID, req)

	return &container.CreateContainerReply{
		BuildID: buildID,
	}, nil

}

// ImageStatus maps to the ImageStatus RPC
func (s *Server) ImageStatus(ctx context.Context, req *container.Build) (*container.ContainerImage, error) {
	// get the build request from the db
	build, err := s.builder.GetBuild(req.Id)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to get build request")
		return nil, errors.Wrap(err, "unable to fetch build")
	}

	return &container.ContainerImage{
		Status: build.Status,
	}, nil
}
