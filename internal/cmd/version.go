package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd(version string) *cobra.Command {

	const (
		use   = "version"
		short = "Print the version"
	)

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}
