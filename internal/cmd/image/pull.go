package image

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/openclosed-dev/docksider/internal/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	all      bool
	quiet    bool
	platform string
}

type puller struct {
	ctx  context.Context
	opts pullOptions
}

func NewPullCmd() *cobra.Command {

	const (
		use   = "pull [OPTIONS] NAME[:TAG|@DIGEST]"
		short = "Download an image from a registry"
	)

	var opts pullOptions

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			puller := puller{cmd.Context(), opts}
			return puller.execute(args[0])
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Suppress verbose output")

	flags.StringVar(&opts.platform, "platform", "",
		"Set platform if server is multi-platform capable")

	return cmd
}

func (p puller) execute(imageName string) error {

	logs.Warn.SetOutput(os.Stderr)
	if !p.opts.quiet {
		logs.Progress.SetOutput(os.Stderr)
	}

	craneOpts := crane.GetOptions(crane.WithContext(p.ctx))

	imageRef, err := name.ParseReference(imageName, craneOpts.Name...)
	if err != nil {
		return err
	}

	platformList, err := parsePlatforms(p.opts.platform)
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "*.tar")
	if err != nil {
		return err
	}
	path := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(path)
	}()

	if err = p.downloadImage(path, imageRef, &craneOpts); err != nil {
		return err
	}

	return p.loadImage(imageName, tempFile, platformList)
}

func (p *puller) downloadImage(path string, imageRef name.Reference, opts *crane.Options) error {

	logs.Progress.Printf("downloading from registry: %s", imageRef.Context().Registry.Name())

	desc, err := remote.Get(imageRef, opts.Remote...)
	if err != nil {
		return err
	}

	image, err := desc.Image()
	if err != nil {
		return err
	}

	if err = crane.Save(image, imageRef.String(), path); err != nil {
		return fmt.Errorf("failed to save %s as tarball: %w", path, err)
	}

	return nil
}

func (p *puller) loadImage(imageName string, reader io.ReadCloser, platformList []ocispec.Platform) error {

	c, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	logs.Progress.Printf("updating image: %s", imageName)

	var options []client.ImageLoadOption
	if p.opts.quiet {
		options = append(options, client.ImageLoadWithQuiet(true))
	}

	if len(platformList) > 0 {
		options = append(options, client.ImageLoadWithPlatforms(platformList...))
	}

	res, err := c.ImageLoad(p.ctx, reader, options...)
	if err != nil {
		return docker.WrapError(err)
	}
	defer res.Close()

	return p.outputResult(res)
}

func (p *puller) outputResult(result client.ImageLoadResult) error {

	var out io.Writer = os.Stderr
	if p.opts.quiet {
		out = io.Discard
	}

	fd, isTerminal := term.GetFdInfo(out)
	return jsonmessage.DisplayJSONMessagesStream(result, out, fd, isTerminal, nil)
}
