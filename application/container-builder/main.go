package main

// go-git is a highly extensible git implementation in pure Go.
// https://github.com/go-git/go-git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// Basic example of how to clone a repository using clone options.
func main() {
	url := "https://github.com/null-channel/eddington.git"
	directory := "./tmp/eddington"

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		panic(err)
	}
	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()

	if err != nil {
		panic(err)
	}

	fmt.Println(ref)

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())

	if err != nil {
		panic(err)
	}

	fmt.Println(commit)
}
