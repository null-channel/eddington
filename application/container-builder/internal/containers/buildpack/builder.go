package image

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/hashicorp/go-hclog"
	"github.com/uptrace/bun"
)

var log = hclog.New(&hclog.LoggerOptions{
	Name:  "container-builder⚡️",
	Level: hclog.LevelFromString("INFO"),
})

type BuildOpt struct {
	RepoURL   string
	ImageName string
	BuildPack string
	Builder   string
}

type BuildRequest struct {
	RepoName string
	BuildID  string
	Status   bool
}
type BuildPackInfo struct {
	BuildPack string
	Builder   string
}

type Builder struct {
	Client *pack.Client
	l      hclog.Logger
	db     *bun.DB
	// used to store built images
	registry string
}

func NewBuilder(db *bun.DB, pack *pack.Client) (*Builder, error) {
	return &Builder{
		l:      log,
		db:     db,
		Client: pack,
	}, nil
}

// DetectLanguage detects the language of the repo
// Would need refactor very soon because github has rate limit
// https://docs.github.com/en/rest/rate-limit/rate-limit?apiVersion=2022-11-28
func (b *Builder) DetectLanguage(repo string) (string, error) {
	repo = strings.TrimPrefix(repo, "https://github.com/")

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/languages", repo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.l.Error("failed to read response body", "error", err)
		return "", err
	}

	var data map[string]int
	err = json.Unmarshal(body, &data)
	if err != nil {
		b.l.Error("failed to unmarshal response body", "error", err)
		return "", err
	}

	for lang := range data {
		return lang, nil
	}

	return "", fmt.Errorf("no language found in the response")
}

func (b *Builder) GetBuildPackInfo(language string) (BuildPackInfo, error) {
	language = strings.ToLower(language)

	switch language {
	case "go":
		return BuildPackInfo{
			BuildPack: "paketo-buildpacks/go",
			Builder:   "paketobuildpacks/builder-jammy-full",
		}, nil
	case "nodejs":
		return BuildPackInfo{
			BuildPack: "paketo-buildpacks/nodejs",
			Builder:   "paketobuildpacks/builder-jammy-full",
		}, nil
	case "python":
		return BuildPackInfo{
			BuildPack: "paketo-buildpacks/python",
			Builder:   "paketobuildpacks/builder-jammy-full",
		}, nil

	default:
		return BuildPackInfo{}, fmt.Errorf("no buildpack found for %s", language)
	}
}

func (b *Builder) CreateImage(opt BuildOpt) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error)

	go func() {
		err := b.Client.Build(ctx, pack.BuildOptions{
			AppPath:    opt.RepoURL,
			Builder:    opt.Builder,
			Image:      opt.ImageName,
			Buildpacks: []string{opt.BuildPack},
		})
		done <- err
	}()
	select {
	case <-ctx.Done():
		// The build was canceled
		return ctx.Err()
	case err := <-done:
		// The build has completed
		if err != nil {
			return err
		}
		// Build completed successfully
		// do something with the image
	}
	return nil
}