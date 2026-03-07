package cmd

import (
	"github.com/openclosed-dev/docksider/internal/cmd/image"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	const (
		use   = "docksider [OPTIONS] COMMAND [ARG...]"
		short = "A Docker-style CLI for pulling and pushing images directly to container registries"
	)

	root := &cobra.Command{
		Use:           use,
		Short:         short,
		SilenceUsage:  true,
		SilenceErrors: false,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		NewLoginCmd(),
		NewPsCmd(),
		image.NewCmd(),
		image.NewImagesCmd(),
		image.NewPullCmd(),
		image.NewPushCmd(),
	)

	return root
}
