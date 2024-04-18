package models

import "errors"

type NullApplication struct {
	ID                     int64                     `bun:",pk,autoincrement" json:"id"`
	Name                   string                    `bun:"name" json:"name"`
	OrgID                  int64                     `bun:"org_id" json:"org_id"`
	ResourceGroupID        int64                     `bun:"resource_group_id" json:"resource_group_id"`
	Namespace              string                    `bun:"namespace" json:"namespace"`
	NullApplicationService []*NullApplicationService `bun:"rel:has-many,join:id=null_application_id" json:"null_application_service"`
}

type BuildType int

const (
	NodeJS BuildType = iota
	Go
	Rust
	Python
	ContainerImage
)

func (b BuildType) String() string {
	return [...]string{"NodeJS", "Go", "Rust", "Python", "ContainerImage"}[b]
}

type NullApplicationService struct {
	ID                int64     `bun:",pk,autoincrement" json:"id"`
	NullApplicationID int64     `bun:"null_application_id" json:"null_application_id"`
	Type              BuildType `bun:"type" json:"type"`
	GitRepo           string    `bun:"gitRepo" json:"git_repo"`
	GitSha            string    `bun:"gitSha" json:"git_sha"`
	Name              string    `bun:"name" json:"name"`
	Image             string    `bun:"image" json:"image"`
	Cpu               string    `bun:"cpu" json:"cpu"`
	Memory            string    `bun:"memory" json:"memory"`
	Storage           string    `bun:"storage" json:"storage"`
	BuildID           string    `bun:"build_id" json:"build_id"`
}

// Validate validates the NullApplicationService
func (n *NullApplicationService) Validate() error {
	if n.Name == "" {
		return errors.New("name is required")
	}
	if n.Type == ContainerImage {
		if n.GitRepo != "" || n.GitSha != "" {
			return errors.New("gitRepo and gitSha are not allowed for container image")
		}
		return nil
	}
	if n.GitRepo == "" {
		return errors.New("gitRepo is required")
	}
	if n.GitSha == "" {
		return errors.New("gitSha is required")
	}
	return nil
}
