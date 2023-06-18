package main

import (
	"context"
	"flag"
	"fmt"
	"time"
	"net"
	"log"

	"github.com/google/go-github/github"
	"google.golang.org/grpc"
	pb "github.com/null-channel/eddington/proto/github-watcher"
)

var (
	port = flag.Int("port", 10001, "github watcher server port")
)

type server struct {
	pb.UnimplementedWatchRepoServiceServer
}

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWatchRepoServiceServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

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
