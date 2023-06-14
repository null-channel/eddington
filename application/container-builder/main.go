package main

import (
	"context"

	"github.com/null-channel/eddington/application/container-builder/internal/containers/dockerfile"
	"github.com/sirupsen/logrus"
)

// go-git is a highly extensible git implementation in pure Go.
// https://github.com/go-git/go-git

// Basic example of how to clone a repository using clone options.
func main() {
	// url := "https://github.com/null-channel/eddington.git"
	// directory := "./tmp/eddington"

	// r, err := git.PlainClone(directory, false, &git.CloneOptions{
	// 	URL:               url,
	// 	RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	// })

	// if err != nil {
	// 	panic(err)
	// }
	// // ... retrieving the branch being pointed by HEAD
	// ref, err := r.Head()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(ref)

	// // ... retrieving the commit object
	// commit, err := r.CommitObject(ref.Hash())

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(commit)
	ctx := context.Background()
	b, err := dockerfile.NewBuilder(ctx)
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
