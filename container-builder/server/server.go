package server

import (
	"context"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/null-channel/eddington/container-builder/proto/container"
	"github.com/vmware-tanzu/kpack-cli/pkg/image"
	"github.com/vmware-tanzu/kpack-cli/pkg/k8s"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Server struct {
	container.UnimplementedContainerServiceServer
	kubeConfig     *rest.Config
	kpackClientSet *k8s.ClientSet
	log            *zap.SugaredLogger
	repoBase       string
	serviceAccount string
}

func NewServer(kubeConfig *rest.Config, logger *zap.SugaredLogger) (*Server, error) {

	kpcs := k8s.DefaultClientSetProvider{}
	cs, err := kpcs.GetClientSet("default")
	if err != nil {
		return nil, err
	}

	return &Server{
		log:            logger,
		kpackClientSet: &cs,
		repoBase:       "nullchannel",
		serviceAccount: "build-container-service-account",
	}, nil
}

// CreateContainer maps to the CreateContainer RPC“
// call using grpcurl:
// grpcurl -plaintext -d '{"repoURL": "your_repo_url", "type": "your_type", "customerID": "your_customer_id"}' localhost:4040 container.ContainerService/CreateContainer
func (s *Server) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.CreateContainerResponse, error) {
	s.log.Info("Create Container Request")

	cusId := strconv.FormatInt(req.CustomerID, 10)

	name := sanitizeRep(req.RepoURL, req.Directory)

	factory := s.getImageFactory(req)
	img, err := factory.MakeImage("kp-test", "default", s.getRepo(name))

	if err != nil {
		slog.Error("unable to create image factory", err.Error())
	}

	img.Labels["customer"] = cusId
	img.Labels["app"] = req.RepoURL + "/" + req.Directory

	img, err = s.kpackClientSet.KpackClient.KpackV1alpha2().Images(s.kpackClientSet.Namespace).Create(ctx, img, metav1.CreateOptions{})
	if err != nil {

		slog.Error("unable to create image ", err.Error())
	}
	return &container.CreateContainerResponse{
		Image:      s.getRepo(name),
		Generation: 1,
	}, nil
}

// ImageStatus maps to the ImageStatus RPC
func (s *Server) BuildStatus(ctx context.Context, req *container.BuildStatusRequest) (*container.BuildStatusResponse, error) {

	name := s.getRepo(sanitizeRep(req.Repo, req.Directory))
	img, err := s.kpackClientSet.KpackClient.KpackV1alpha2().Images(s.kpackClientSet.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	//TODO: understand how to tell if there is a build happening right now.

	return &container.BuildStatusResponse{
		ImageName: img.Status.LatestImage,
	}, nil
}

func (s *Server) getImageFactory(req *container.CreateContainerRequest) *image.Factory {

	builder := selectBuilder(req)

	return &image.Factory{
		Builder:        builder,
		GitRepo:        req.RepoURL,
		GitRevision:    req.Rev,
		ServiceAccount: s.serviceAccount,
		SubPath:        &req.Directory,
	}
}

func (s *Server) getRepo(gitRepo string) string {
	return s.repoBase + "/" + gitRepo
}

func (s *Server) getImageName(cusId, gitrepo string) string {
	return cusId + "-" + gitrepo
}

func sanitizeRep(repo, dir string) string {
	return strings.Replace((repo + "/" + dir), "/", "-", -1)
}

// selectBuilder TODO: select a builder to build the project. not sure if we can just have one that auto selects or if we will have to have multiples.
func selectBuilder(req *container.CreateContainerRequest) string {
	return "main-builder"
}
