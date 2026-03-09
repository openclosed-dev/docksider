package docker

import (
	"fmt"
	"os"
	"runtime"

	"github.com/moby/moby/client"
)

func NewClient() (*client.Client, error) {

	if runtime.GOOS == "windows" {
		value := os.Getenv("DOCKER_HOST")
		if value == "" {
			return nil, fmt.Errorf(`The environment variable DOCKER_HOST is not defined.
It must have a value similar to 'tcp://<host>:<port>'.`)
		}
	}

	return client.New(client.WithHostFromEnv())
}

func WrapError(err error) error {
	if client.IsErrConnectionFailed(err) {
		return fmt.Errorf(`failed to connect to the Docker daemon: %w\n
Hint: If the daemon is installed on a WSL distribution, start the instance first.
`, err)
	}
	return err
}
