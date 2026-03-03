package cmd

import "github.com/spf13/cobra"

// Hidden command invoked by `az acr login`.
func NewPsCmd() *cobra.Command {

	const (
		use   = "ps [OPTIONS]"
		short = "List containers"
	)

	cmd := &cobra.Command{
		Use:    use,
		Short:  short,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
