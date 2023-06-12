package dockerfile

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/client"
	"github.com/sirupsen/logrus"
)

type BuildOpt struct {
	Dockerfile string
	ImageName  string
	Tag        string
}

type Builder struct {
	client *client.Client
}

func NewBuilder(ctx context.Context) (*Builder, error) {
	client, err := client.New(ctx, "", client.WithFailFast())
	if err != nil {
		logrus.Panic("unable to create buildkit client error: ", err.Error())
		return nil, err
	}
	return &Builder{
		client: client,
	}, nil

}

func (b *Builder) Build(opts BuildOpt) error {
	ctx := context.Background()
	solvOpts := b.createSolveOpt(opts.ImageName, ".", opts.Dockerfile)

	// prob print out what's going on here
	status := make(chan *client.SolveStatus)

	b.client.Solve(ctx, nil, solvOpts, status)
	return nil

}

func (b *Builder) createSolveOpt(imageName string, buildContext string, dockerfile string) client.SolveOpt {
	return client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: "image",
				Attrs: map[string]string{
					"name": imageName,
				},
				Output: func(m map[string]string) (io.WriteCloser, error) {
					return os.Stdout, nil
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
