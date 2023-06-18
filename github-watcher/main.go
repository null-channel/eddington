package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/go-github/github"
	"github.com/nats-io/nats.go"
	pb "github.com/null-channel/eddington/proto/github-watcher"
	"google.golang.org/grpc"
)

// var (
// 	// port = flag.Int("port", 10001, "github watcher server port")

// )

// request -> Watch(xyz/project:main)
// User [watchedRepos(xyz/project:main)]
// for rep in watched repos GetUpdates?

type WatchRepoServer struct {
	pb.UnimplementedWatchRepoServiceServer
	nc *nats.Conn
	// todo: github client
	gh *github.Client
}

func (s *WatchRepoServer) WatchRepo(ctx context.Context, req *pb.WatchRepoRequest) (*pb.WatchRepoResponse, error) {

	// nc, ncErr := nats.Connect("localhost:4222")
	defer s.nc.Drain()

	// if ncErr != nil {
	// 	log.Fatal(ncErr)
	// }
	defer s.nc.Close()

	res := s.nc.Status()
	log.Printf("Status: %+v", res)

	pollCommits(req.Owner, req.Repository, req.Branch, req.SHA, s.gh, s.nc)

	return &pb.WatchRepoResponse{IsUpdated: true}
}

//

func pollCommits(owner string, repo string, branch string, lastCommitSHA string, client *github.Client, nc *nats.Conn) {

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

		if err := nc.Publish("container-builder", []byte("This branch has been updated")); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// client := github.NewClient(nil)

	// nc, ncErr := nats.Connect("localhost:4222")
	// defer nc.Drain()

	// if ncErr != nil {
	// 	log.Fatal(ncErr)
	// }
	// defer nc.Close()

	// res := nc.Status()
	// log.Printf("Status: %+v", res)

	// // owner := "null-channel"
	// // repo := "eddington"
	// // branch := "main"
	// // lastCommitSHA := "02514afee6fc9ada5225a4f9670ebc4e627d6e24"

	// usePolling := flag.Bool("poll", false, "Continuously poll github repo every five minutes")

	// flag.Parse()

	// if *usePolling {
	// 	for {
	// 		pollCommits(owner, repo, branch, lastCommitSHA, client, nc)

	// 		time.Sleep(5 * time.Minute)
	// 	}
	// } else {
	// 	pollCommits(owner, repo, branch, lastCommitSHA, client, nc)
	// }

	// lis
	lis, err := net.Listen("tcp", "10001")
	if err != nil {
		log.Panicf("could not open shop %s", err.Error())
	}

	// rpc server
	server := grpc.NewServer()
	nc, _ := nats.Connect(nats.DefaultURL)
	gh := github.NewClient(nil)

	pb.RegisterWatchRepoServiceServer(server, &WatchRepoServer{
		nc: nc,
		gh: gh,
	})

	err = server.Serve(lis)
	if err != nil {
		log.Panicf("oops couldn't open up shop %s", err.Error())
	}

}
