package docker

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/moby/moby/client"
)

func NewClient() (*client.Client, error) {

	host, ok := os.LookupEnv("DOCKER_HOST")
	if ok {
		err := ValidateHost(host)
		if err != nil {
			return nil, fmt.Errorf("DOCKER_HOST has an invalid value '%s': %w", host, err)
		}
	} else {
		host = client.DefaultDockerHost
	}

	return NewClientForHost(host)
}

func NewClientForHost(host string) (*client.Client, error) {
	return client.New(client.WithHost(host))
}

var protocolsForUnix = []string{"tcp", "ssh", "unix"}
var protocolsForWindows = []string{"tcp", "ssh", "npipe"}

func ValidateHost(value string) error {

	if len(strings.TrimSpace(value)) == 0 {
		return errors.New("value is blank")
	}

	scheme, remaining, ok := strings.Cut(value, "://")
	if !ok || scheme == "" {
		return errors.New("protocol is missing in the URL")
	}

	switch scheme {
	case "tcp":
		if err := validateTcpHost(value); err != nil {
			return err
		}
	case "unix":
		if err := validateUnixHost(remaining); err != nil {
			return err
		}
	case "ssh", "npipe":
	default:
		return fmt.Errorf(
			"unsupported protocol '%s'; "+
				"the protocol must be one of [%s] for Unix and [%s] for Windows; ",
			scheme,
			strings.Join(protocolsForUnix, ", "),
			strings.Join(protocolsForWindows, ", "))
	}

	return nil
}

func WrapError(err error) error {
	if client.IsErrConnectionFailed(err) {
		return fmt.Errorf(
			`failed to connect to the Docker daemon: %w
Hint: If the daemon is installed on a WSL distribution, start the distribution first.`,
			err)
	}
	return err
}

func validateTcpHost(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("not a URL: %w", err)
	}
	if parsed.Hostname() == "" {
		return errors.New("hostname is missing")
	}
	if parsed.Port() == "" {
		return errors.New("port number is missing")
	}
	return nil
}

func validateUnixHost(path string) error {
	if runtime.GOOS == "windows" {
		return errors.New("'unix' protocol is not supported in Windows")
	}
	if len(strings.TrimSpace(path)) == 0 {
		return errors.New("path is blank")
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path is not absolute: '%s'", path)
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file does not exist at '%s'", path)
	}
	return nil
}
