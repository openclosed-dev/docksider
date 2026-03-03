package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	username      string
	password      string
	passwordStdin bool
}

func NewLoginCmd() *cobra.Command {

	const (
		use   = "login [OPTIONS] [SERVER]"
		short = "Authenticate to a registry"
		long  = "Authenticate to a registry.\nDefaults to Docker Hub if no server is specified."
	)

	var opts loginOptions

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(args[0], opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.username, "username", "u", "", "Username")
	flags.StringVarP(&opts.password, "password", "p", "", "Password or Personal Access Token (PAT)")
	flags.BoolVar(&opts.passwordStdin, "password-stdin", false, "Take the Password or Personal Access Token (PAT) from stdin")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagsOneRequired("password", "password-stdin")
	cmd.MarkFlagsMutuallyExclusive("password", "password-stdin")

	return cmd
}

func login(server string, opts loginOptions) error {

	logs.Warn.SetOutput(os.Stderr)
	logs.Progress.SetOutput(os.Stderr)

	server, err := getCanonicalRegistryName(server)
	if err != nil {
		return nil
	}

	if opts.passwordStdin {
		password, err := readPasswordFromStdin()
		if err != nil {
			return err
		}
		opts.password = password
	}

	cf, err := config.Load(os.Getenv("DOCKER_CONFIG"))
	if err != nil {
		return err
	}

	creds := cf.GetCredentialsStore(server)
	if server == name.DefaultRegistry {
		server = authn.DefaultAuthKey
	}

	auth := types.AuthConfig{
		ServerAddress: server,
		Username:      opts.username,
		Password:      opts.password,
	}

	if err := creds.Store(auth); err != nil {
		return err
	}

	if err := cf.Save(); err != nil {
		return err
	}

	logs.Progress.Printf("logged in via %s", cf.Filename)
	return nil
}

func getCanonicalRegistryName(server string) (string, error) {
	reg, err := name.NewRegistry(server)
	if err != nil {
		return "", err
	}
	return reg.Name(), nil
}

func readPasswordFromStdin() (string, error) {
	contents, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	password := string(contents)
	password = strings.TrimSuffix(password, "\n")
	password = strings.TrimSuffix(password, "\r")
	return password, nil
}
