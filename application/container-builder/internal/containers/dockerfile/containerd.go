package dockerfile

import (
	"context"

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

}
