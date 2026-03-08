package cmd

import (
	"strings"

	"github.com/openclosed-dev/docksider/internal/cmd/image"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {

	const (
		use   = "docksider [OPTIONS] COMMAND [ARG...]"
		short = "A Docker-style CLI for pulling and pushing images directly to container registries"
	)

	version = strings.TrimSpace(version)

	root := &cobra.Command{
		Use:           use,
		Short:         short,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: false,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		DisableFlagsInUseLine: true,
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		NewLoginCmd(),
		NewPsCmd(),
		image.NewCmd(),
		image.NewImagesCmd(),
		image.NewPullCmd(),
		image.NewPushCmd(),
		NewVersionCmd(version),
	)

	return root
}
