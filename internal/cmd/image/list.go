package image

import (
	"context"
	"os"

	"github.com/docker/cli/cli/command/formatter"
	flagsHelper "github.com/docker/cli/cli/flags"
	"github.com/docker/cli/opts"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/client"
	"github.com/openclosed-dev/docksider/internal/docker"
	"github.com/spf13/cobra"
)

type listOptions struct {
	quiet   bool
	all     bool
	noTrunc bool
	digests bool
	tree    bool
	format  string
	filter  opts.FilterOpt
}

type lister struct {
	ctx  context.Context
	opts listOptions
}

func NewListCmd() *cobra.Command {

	const (
		use   = "ls [OPTIONS] [REPOSITORY[:TAG]]"
		short = "List images"
	)

	var opts = listOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Aliases: []string{"list"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := lister{
				ctx:  cmd.Context(),
				opts: opts,
			}
			return l.execute(args)
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()

	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show image IDs")
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all images (default hides intermediate and dangling images)")
	flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")
	flags.BoolVar(&opts.digests, "digests", false, "Show digests")
	flags.StringVar(&opts.format, "format", "", flagsHelper.FormatHelp)
	flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")

	return cmd
}

func NewImagesCmd() *cobra.Command {
	const (
		use = "images [OPTIONS] [REPOSITORY[:TAG]]"
	)

	cmd := NewListCmd()
	cmd.Use = use
	cmd.Aliases = []string{"image ls", "image list"}

	return cmd
}

func (l *lister) execute(args []string) error {

	filters := l.opts.filter.Value()
	if len(args) > 0 {
		filters.Add("reference", args[0])
	}

	c, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	listOpts := client.ImageListOptions{
		All:       l.opts.all,
		Filters:   filters,
		Manifests: false,
	}

	res, err := c.ImageList(l.ctx, listOpts)
	if err != nil {
		return err
	}

	return l.outputResult(res.Items)
}

func (l *lister) outputResult(images []image.Summary) error {

	format := formatter.TableFormatKey

	imageCtx := formatter.ImageContext{
		Context: formatter.Context{
			Output: os.Stdout,
			Format: formatter.NewImageFormat(format, l.opts.quiet, l.opts.digests),
			Trunc:  !l.opts.noTrunc,
		},
		Digest: l.opts.digests,
	}

	return formatter.ImageWrite(imageCtx, images)
}
