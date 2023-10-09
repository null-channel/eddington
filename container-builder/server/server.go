package server

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/null-channel/eddington/proto/container"
	"github.com/vmware-tanzu/kpack-cli/pkg/image"
	"github.com/vmware-tanzu/kpack-cli/pkg/k8s"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Server struct {
	container.UnimplementedContainerServiceServer
	kubeConfig *rest.Config
	log        *zap.SugaredLogger
}

func NewServer(kubeConfig *rest.Config, logger *zap.SugaredLogger) (*Server, error) {

	return &Server{
		log: logger,
	}, nil
}

// CreateContainer maps to the CreateContainer RPC“
// call using grpcurl:
// grpcurl -plaintext -d '{"repoURL": "your_repo_url", "type": "your_type", "customerID": "your_customer_id"}' localhost:4040 container.ContainerService/CreateContainer
func (s *Server) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.CreateContainerResponse, error) {
	s.log.Info("Create Container Request")

	factory := image.Factory{
		Builder:        "my-builder",
		GitRepo:        req.RepoURL,
		ServiceAccount: "tutorial-service-account",
		GitRevision:    "main",
	}

	img, err := factory.MakeImage("kp-test", "default", "nullchannel/nc-test-random")

	if err != nil {
		slog.Error("unable to create image factory", err.Error())

	}

	kpcs := k8s.DefaultClientSetProvider{}
	cs, err := kpcs.GetClientSet("default")
	img, err = cs.KpackClient.KpackV1alpha2().Images(cs.Namespace).Create(ctx, img, metav1.CreateOptions{})
	if err != nil {

		slog.Error("unable to create image ", err.Error())
	}
	return &container.CreateContainerResponse{
		BuildID: "1234",
	}, nil
}

// ImageStatus maps to the ImageStatus RPC
func (s *Server) BuildStatus(ctx context.Context, req *container.BuildStatusRequest) (*container.BuildStatusResponse, error) {

	return nil, nil

}
