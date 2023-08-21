package models

import "github.com/null-channel/eddington/proto/container"

type Build struct {
	CustomerID    string                    `bun:"customer_id,notnull"`
	RepoName      string                    `bun:"name,notnull"`
	BuildID       string                    `bun:"build_id,pk"`
	Status        container.ContainerStatus `bun:"status,default:0"`
	StatusMessage string                    `bun:"status_message"`
}
