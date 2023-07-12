package main

import (
	"context"

	"github.com/null-channel/eddington/application/container-builder/internal/containers/dockerfile"
	"github.com/sirupsen/logrus"
)

var (
	builkitdaddr = "tcp://0.0.0.0:4000"
)

func main() {

	ctx := context.Background()
	b, err := dockerfile.NewBuilder(ctx, builkitdaddr)
	if err != nil {
		logrus.Panic("unable to create builder error: ", err.Error())
	}
	opts := dockerfile.BuildOpt{
		Dockerfile: "./Dockerfile",
		ImageName:  "eddington",
	}
	err = b.Build(opts)
	if err != nil {
		logrus.Panic("unable to build image error: ", err.Error())
	}

}
