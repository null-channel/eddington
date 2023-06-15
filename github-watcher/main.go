package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

func pollCommits(owner string, repo string, branch string, lastCommitSHA string, client *github.Client) {

	ctx := context.Background()

	opt := &github.CommitsListOptions{
		SHA: branch,
	}

	commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	latestCommitSHA := *commits[0].SHA
	if latestCommitSHA == lastCommitSHA {
		fmt.Println("No updates to this branch")
	} else {
		fmt.Println("This branch has been updated")
	}

}

func main() {
	client := github.NewClient(nil)

	owner := "null-channel"
	repo := "eddington"
	branch := "main"
	lastCommitSHA := "02514afee6fc9ada5225a4f9670ebc4e627d6e24"

	usePolling := flag.Bool("poll", false, "Continuously poll github repo every five minutes")

	flag.Parse()

	if *usePolling {
		for {
			pollCommits(owner, repo, branch, lastCommitSHA, client)

			time.Sleep(5 * time.Minute)
		}
	} else {
		pollCommits(owner, repo, branch, lastCommitSHA, client)
	}
}
