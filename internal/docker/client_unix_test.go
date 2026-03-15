//go:build !windows

package docker_test

import (
	"strings"
	"testing"

	"github.com/openclosed-dev/docksider/internal/docker"
)

func TestValidHostOnUnix(t *testing.T) {

	cases := []struct {
		name string
		host string
	}{
		{"unix", "unix:///var/run/docker.sock"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := docker.ValidateHost(c.host)
			if err != nil {
				t.Errorf("failed to validate: %s", err)
			}
		})
	}
}

func TestInvalidHostOnUnix(t *testing.T) {

	cases := []struct {
		name    string
		host    string
		message string
	}{
		{"unix", "unix://", "path is blank"},
		{"unix", "unix://localhost/var/run/docker.sock", "path is not absolute"},
		{"unix", "unix:///nonexistent", "file does not exist"},
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
