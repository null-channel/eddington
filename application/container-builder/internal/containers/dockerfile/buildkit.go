package dockerfile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/continuity"
	"github.com/google/uuid"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/cmd/buildctl/build"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/util/bklog"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BuildOpt struct {
	Dockerfile string
	ImageName  string
	Tag        string
}

type Builder struct {
	Client *client.Client
	// TODO: add nats connection
	// nc     *nats.Conn
}

func NewBuilder(ctx context.Context, builkitdaddr string) (*Builder, error) {
	client, err := client.New(ctx, builkitdaddr, client.WithFailFast())
	if err != nil {
		logrus.Panic("unable to create buildkit client error: ", err.Error())
		return nil, err
	}
	return &Builder{
		Client: client,
	}, nil

}

func (b *Builder) Build(opts BuildOpt) error {

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	exports, err := build.ParseOutput([]string{"type=image,name=" + opts.ImageName})
	if err != nil {
		return errors.New("unable to parse output error: " + err.Error())
	}

	ref := uuid.New().String()

	solveOpt := client.SolveOpt{
		Exports:  exports,
		Frontend: "dockerfile.v0",
		Ref:      ref,
		LocalDirs: map[string]string{
			"context":    ".",
			"dockerfile": filepath.Dir(opts.Dockerfile),
		},
		FrontendAttrs: map[string]string{
			"filename": opts.Dockerfile,
		},
	}
	pw, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, "auto")
	if err != nil {
		return err
	}
	mw := progresswriter.NewMultiWriter(pw)
	var writers []progresswriter.Writer

	var subMetadata map[string][]byte

	eg.Go(func() error {
		defer func() {
			for _, w := range writers {
				close(w.Status())
			}
		}()

		sreq := gateway.SolveRequest{
			Frontend:    solveOpt.Frontend,
			FrontendOpt: solveOpt.FrontendAttrs,
		}

		resp, err := b.Client.Build(ctx, solveOpt, "buildctl", func(ctx context.Context, c gateway.Client) (*gateway.Result, error) {
			_, isSubRequest := sreq.FrontendOpt["requestid"]
			if isSubRequest {
				if _, ok := sreq.FrontendOpt["frontend.caps"]; !ok {
					sreq.FrontendOpt["frontend.caps"] = "moby.buildkit.frontend.subrequests"
				}
			}
			res, err := c.Solve(ctx, sreq)
			if err != nil {
				return nil, err
			}
			if isSubRequest && res != nil {
				subMetadata = res.Metadata
			}
			return res, err
		}, progresswriter.ResetTime(mw.WithPrefix("", false)).Status())
		if err != nil {
			return err
		}
		for k, v := range resp.ExporterResponse {
			bklog.G(ctx).Debugf("exporter response: %s=%s", k, v)
		}

		metadataFile := "./metadata.txt"
		if metadataFile != "" && resp.ExporterResponse != nil {
			if err := writeMetadataFile(metadataFile, resp.ExporterResponse); err != nil {
				return err
			}
		}

		return nil
	})

	eg.Go(func() error {
		<-pw.Done()
		return pw.Err()
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if txt, ok := subMetadata["result.txt"]; ok {
		fmt.Print(string(txt))
	} else {
		for k, v := range subMetadata {
			if strings.HasPrefix(k, "result.") {
				fmt.Printf("%s\n%s\n", k, v)
			}
		}
	}
	return nil
}

func writeMetadataFile(filename string, exporterResponse map[string]string) error {
	var err error
	out := make(map[string]interface{})
	for k, v := range exporterResponse {
		dt, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			out[k] = v
			continue
		}
		var raw map[string]interface{}
		if err = json.Unmarshal(dt, &raw); err != nil || len(raw) == 0 {
			out[k] = v
			continue
		}
		out[k] = json.RawMessage(dt)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return continuity.AtomicWriteFile(filename, b, 0666)
}
