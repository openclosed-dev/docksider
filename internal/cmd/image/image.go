package image

import (
	"fmt"

	"github.com/containerd/platforms"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {

	const (
		use   = "image"
		short = "Manage images"
	)

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
	}

	cmd.AddCommand(
		NewListCmd(),
		NewPullCmd(),
		NewPushCmd(),
	)

	return cmd
}

func parsePlatforms(platform string) ([]ocispec.Platform, error) {

	list := []ocispec.Platform{}
	if platform != "" {
		parsed, err := platforms.Parse(platform)
		if err != nil {
			return nil, fmt.Errorf("invalid platform: %w", err)
		}
		list = append(list, parsed)
	}

	return list, nil
}
