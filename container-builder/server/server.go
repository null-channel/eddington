package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"

	"strings"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/null-channel/eddington/container-builder/git"
	image "github.com/null-channel/eddington/container-builder/internal/containers/buildpack"
	"github.com/null-channel/eddington/container-builder/models"
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
	// directory where cloned repos would be stored
	repoDirs string
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

	registry := os.Getenv("REGISTRY_URL")
	if registry == "" {
		return nil, errors.New("REGISTRY_URL is not set")
	}

	builder, err := image.NewBuilder(db, packClient, registry)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create builder")
	}
	repoDirs := GetRepoDirOrDefault()
	// create the repo directory if it doesn't exist
	if _, err := os.Stat(repoDirs); os.IsNotExist(err) {
		err = os.Mkdir(repoDirs, 0755)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create repo directory")
		}
	}

	l := zerolog.New(os.Stdout).With().Timestamp().Str("component", "server").Logger()
	return &Server{
		builder:  builder,
		repoDirs: repoDirs,
		log:      l,
	}, nil

}

// CreateContainer maps to the CreateContainer RPC``
// call using grpcurl:
// grpcurl -plaintext -d '{"repoURL": "your_repo_url", "type": "your_type", "customerID": "your_customer_id"}' localhost:4040 container.ContainerService/CreateContainer

func (s *Server) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.CreateContainerResponse, error) {
	s.log.Info().Msg("processing build for repo " + req.RepoURL)
	buildID := uuid.New().String()
	// create a build request in the db

	// validate the repo url
	isValid, _ := git.ValidateGitHubURL(req.RepoURL)

	if !isValid {
		return nil, errors.New("invalid repo url")
	}

	repo := strings.TrimPrefix(req.RepoURL, "https://github.com/")
	err := s.builder.NewBuild(req.CustomerID, repo, buildID)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to create build request")
		return nil, errors.Wrap(err, "unable to create build , please try again")
	}

	dir := path.Join(s.repoDirs, fmt.Sprintf("%s-%s", buildID, req.CustomerID))
	s.log.Debug().Str("dir", dir).Msg("cloning repo")

	// start build in a goroutine with the buildID and request
	go func(buildID string, req *container.CreateContainerRequest) {
		//TODO: Validate git repo
		_, err = git.Clone(req.RepoURL, dir)
		if err != nil {
			s.log.Error().Err(err).Msg("unable to clone repo")
			// Failed to clone the repo
			s.builder.UpdateBuildRequest(buildID, container.ContainerStatus_FAILED, "unable to clone repo")
			return
		}
		buildInfo, err := s.builder.GetBuildPackInfo(req.Type.String())
		if err != nil {
			s.log.Error().Err(err).Msg("unable to get buildpack info")
			s.builder.UpdateBuildRequest(buildID, container.ContainerStatus_FAILED, "unable to retreive buildpack info")
			return
		}

		// add repo name in the format of registry/customerID-repoName
		// imagName := fmt.Sprintf("%s/%s-%s", s.builder.Registry, req.CustomerID, repo)
		buildPath := path.Join(dir, req.Directory)

		opts := image.BuildOpt{
			BuildID:   buildID,
			RepoURL:   req.RepoURL,
			Path:      buildPath,
			ImageName: fmt.Sprintf("%s/%s-%s", s.builder.Registry, "nc-test", "random"),
			BuildPack: buildInfo.BuildPack,
			Builder:   buildInfo.Builder,
		}

		err = s.builder.CreateImage(opts)
		if err != nil {
			s.log.Error().Err(err).Msg("Build failed")
			s.builder.UpdateBuildRequest(buildID, container.ContainerStatus_FAILED, err.Error())
			return
		}

	}(buildID, req)

	return &container.CreateContainerResponse{
		BuildID: buildID,
	}, nil

}

// ImageStatus maps to the ImageStatus RPC
func (s *Server) BuildStatus(ctx context.Context, req *container.BuildStatusRequest) (*container.BuildStatusResponse, error) {
	// get the build request from the db
	build, err := s.builder.GetBuild(req.Id)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to get build request")
		return nil, errors.Wrap(err, "unable to fetch build")
	}

	return &container.BuildStatusResponse{
		ImageName: build.RepoName,
		Status:    build.Status,
	}, nil
}

func GetRepoDirOrDefault() string {
	env := os.Getenv("ENVIROMENT")
	if env == "dev" || env == "" {
		return "./tmp"
	} else {
		return "/tmp"
	}
}
