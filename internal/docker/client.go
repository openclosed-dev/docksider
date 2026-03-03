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

	return client.NewClientWithOpts(client.WithHostFromEnv())
}
