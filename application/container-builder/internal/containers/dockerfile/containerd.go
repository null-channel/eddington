package dockerfile

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/buildkit/client"
	dockerfile "github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BuildOpt struct {
	Dockerfile string
	ImageName  string
	Tag        string
}

type Builder struct {
	Client *client.Client
	// TODO: add nats connection
	// nc     *nats.Conn
}

func NewBuilder(ctx context.Context) (*Builder, error) {
	client, err := client.New(ctx, "tcp://0.0.0.0:4000", client.WithFailFast())
	if err != nil {
		logrus.Panic("unable to create buildkit client error: ", err.Error())
		return nil, err
	}
	return &Builder{
		Client: client,
	}, nil

}

func (b *Builder) Build(opts BuildOpt) error {

	ctx := context.Background()
	pipeR, pipeW := io.Pipe()
	eg, ctx := errgroup.WithContext(ctx)

	solvOpts := b.createSolveOpt(opts.ImageName, ".", opts.Dockerfile, pipeW)

	// prob print out what's going on here
	status := make(chan *client.SolveStatus)

	eg.Go(func() error {
		var err error
		_, err = b.Client.Build(ctx, solvOpts, "", dockerfile.Build, status)
		if err != nil {
			logrus.Panic("unable to build image error: ", err.Error())
			return err
		}
		return nil
	})

	eg.Go(func() error {
		if err := loadDockerTar(pipeR); err != nil {
			return err
		}
		return pipeR.Close()
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	logrus.Infof("Loaded the image %q to Docker.")
	return nil

}

func (b *Builder) createSolveOpt(imageName string, buildContext string, dockerfile string, w io.WriteCloser) client.SolveOpt {

	return client.SolveOpt{
		Exports: []client.ExportEntry{
			{

				Type: "image",
				Attrs: map[string]string{
					"name": imageName,
				},

				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		LocalDirs: map[string]string{
			"context":    buildContext,
			"dockerfile": filepath.Dir(dockerfile),
		},
		Frontend: "dockerfile.v0",
		FrontendAttrs: map[string]string{
			"filename": filepath.Base(dockerfile),
		},
	}
}

func loadDockerTar(r io.Reader) error {

	cmd := exec.Command("docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
