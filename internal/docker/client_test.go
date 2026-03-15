package docker_test

import (
	"strings"
	"testing"

	"github.com/openclosed-dev/docksider/internal/docker"
)

func TestValidHost(t *testing.T) {

	cases := []struct {
		name string
		host string
	}{
		{"tcp", "tcp://192.168.0.100:2375"},
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

func TestInvalidHost(t *testing.T) {

	cases := []struct {
		name    string
		host    string
		message string
	}{
		{"empty", "", "value is blank"},
		{"blank", " ", "value is blank"},
		{"missing protocol", "192.168.0.100:2375", "protocol is missing"},
		{"unsupported protocol", "http://localhost", "unsupported protocol"},
		{"missing host", "tcp://", "hostname is missing"},
		{"missing port", "tcp://192.168.0.100", "port number is missing"},
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
