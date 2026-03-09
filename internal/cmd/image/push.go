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
	"github.com/openclosed-dev/docksider/internal/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	all      bool
	quiet    bool
	platform string
}

type pusher struct {
	ctx  context.Context
	opts pushOptions
}

func NewPushCmd() *cobra.Command {

	const (
		use   = "push [OPTIONS] NAME[:TAG]"
		short = "Upload an image to a registry"
	)

	var opts pushOptions

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pusher := pusher{cmd.Context(), opts}
			return pusher.execute(args[0])
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Suppress verbose output")

	flags.StringVar(&opts.platform, "platform", "",
		`Push a platform-specific manifest as a single-platform image to the registry.
Image index won't be pushed, meaning that other manifests, including attestations won't be preserved.
'os[/arch[/variant]]': Explicit platform (eg. linux/amd64)`)

	return cmd
}

func (p *pusher) execute(imageName string) error {

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
	defer os.Remove(path)

	if err := p.saveImage(imageName, tempFile, platformList); err != nil {
		return err
	}

	return p.uploadImage(path, imageRef, &craneOpts)
}

func (p *pusher) saveImage(imageName string, writer io.WriteCloser, platformList []ocispec.Platform) error {

	defer writer.Close()

	c, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	images := []string{imageName}

	var saveOpts []client.ImageSaveOption
	if len(platformList) > 0 {
		saveOpts = append(saveOpts, client.ImageSaveWithPlatforms(platformList...))
	}

	logs.Progress.Printf("preparing image: %s", imageName)

	reader, err := c.ImageSave(p.ctx, images, saveOpts...)
	if err != nil {
		return docker.WrapError(err)
	}
	defer reader.Close()

	_, err = io.Copy(writer, reader)
	return err
}

func (p *pusher) uploadImage(path string, imageRef name.Reference, opts *crane.Options) error {

	image, err := crane.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load %s as tarball: %w", path, err)
	}

	logs.Progress.Printf("uploading to registry: %s", imageRef.Context().Registry.Name())

	if err := remote.Write(imageRef, image, opts.Remote...); err != nil {
		return err
	}

	if logs.Enabled(logs.Progress) {
		h, err := image.Digest()
		if err != nil {
			return err
		}
		digest := imageRef.Context().Digest(h.String())
		fmt.Fprintln(os.Stdout, digest)
	}

	return nil
}
