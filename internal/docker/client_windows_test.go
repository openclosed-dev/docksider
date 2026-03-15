//go:build windows

package docker_test

import (
	"strings"
	"testing"

	"github.com/openclosed-dev/docksider/internal/docker"
)

func TestInvalidHostOnWindows(t *testing.T) {

	cases := []struct {
		name    string
		host    string
		message string
	}{
		{"unix", "unix:///var/run/docker.sock", "'unix' protocol is not supported in Windows"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := docker.ValidateHost(c.host)
			if err == nil {
				t.Error("failed to detect problem")
			}
			message := err.Error()
			if !strings.HasPrefix(message, c.message) {
				t.Errorf("unexpected error message: %s", message)
			}
		})
	}
}
