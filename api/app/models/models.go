package models

import "errors"

type NullApplication struct {
	ID                     int64                     `bun:",pk,autoincrement"`
	Name                   string                    `bun:"name"`
	OrgID                  int64                     `bun:"org_id"`
	ResourceGroupID        int64                     `bun:"resource_group_id"`
	Namespace              string                    `bun:"namespace"`
	NullApplicationService []*NullApplicationService `bun:"rel:has-many,join:id=null_application_id"`
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
	ID                int64     `bun:",pk,autoincrement"`
	NullApplicationID int64     `bun:"null_application_id"`
	Type              BuildType `bun:"type"`
	GitRepo           string    `bun:"gitRepo"`
	GitSha            string    `bun:"gitSha"`
	Name              string    `bun:"name"`
	Image             string    `bun:"image"`
	Cpu               string    `bun:"cpu"`
	Memory            string    `bun:"memory"`
	Storage           string    `bun:"storage"`
	BuildID           string    `bun:"build_id"`
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
