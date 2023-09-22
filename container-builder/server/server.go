package server

import (
	"context"
	"fmt"
	"os"
	"path"

	"strings"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/null-channel/eddington/container-builder/git"
	image "github.com/null-channel/eddington/container-builder/internal/containers/buildpack"
	"github.com/null-channel/eddington/container-builder/internal/containers/templates"
	"github.com/null-channel/eddington/proto/container"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type Server struct {
	container.UnimplementedContainerServiceServer
	builder *image.Builder
	log     zerolog.Logger
	// directory where cloned repos would be stored
	repoDirs string
}

func NewServer(db *bun.DB, builder *image.Builder, logger *zerolog.Logger) (*Server, error) {

	repoDirs := GetRepoDirOrDefault()
	// create the repo directory if it doesn't exist
	if _, err := os.Stat(repoDirs); os.IsNotExist(err) {
		err = os.Mkdir(repoDirs, 0755)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create repo directory")
		}
	}

	return &Server{
		builder:  builder,
		repoDirs: repoDirs,
		log:      *logger,
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
	err := git.ValidateGitHubURL(req.RepoURL)

	if err != nil {
		return nil, err
	}

	repo := strings.TrimPrefix(req.RepoURL, "https://github.com/")
	err = s.builder.NewBuild(req.CustomerID, repo, buildID)
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
		defer func() {
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Println("Failed to delete dir: " + err.Error())
			}
		}()
		if err != nil {
			s.log.Error().Err(err).Msg("unable to clone repo")
			// Failed to clone the repo
			s.builder.UpdateBuildRequest(buildID, container.ContainerStatus_FAILED, "unable to clone repo")
			return
		}

		s.log.Debug().Msg("Cloned Repo")
		// add repo name in the format of registry/customerID-repoName
		// imagName := fmt.Sprintf("%s/%s-%s", s.builder.Registry, req.CustomerID, repo)
		buildPath := path.Join(dir, req.Directory)

		buildInfo, err := s.builder.GetBuildPackInfo(req.Type.String())
		if err != nil {
			s.log.Error().Err(err).Msg("unable to get buildpack info")
			s.builder.UpdateBuildRequest(buildID, container.ContainerStatus_FAILED, "unable to retreive buildpack info")
			return
		}

		s.log.Debug().Msg("Obtained Build Pack Info")
		language := strings.ToLower(req.Type.String())

		//TODO: Do we need to have a "public" folder and the nginx config file outside that?
		if language == "static-web" {
			//write nginx config to filesystem
			err = os.WriteFile(buildPath, []byte(templates.NginxConf), 0777)
		}

		opts := image.BuildOpt{
			BuildID:   buildID,
			RepoURL:   req.RepoURL,
			Path:      buildPath,
			ImageName: fmt.Sprintf("%s/%s-%s", s.builder.Registry, "nc-test", "random"),
			BuildPack: buildInfo.BuildPack,
			Builder:   buildInfo.Builder,
		}
		client, err := pack.NewClient()
		if err != nil {
			s.log.Err(err).Msgf("unable to create pack client for build: %s", buildID)
		}
		err = s.builder.CreateImage(opts, client)
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
