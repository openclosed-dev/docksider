package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/client"
	"github.com/openclosed-dev/docksider/internal/docker"
	"github.com/spf13/cobra"
)

type diagnoser struct {
	ctx              context.Context
	out              io.Writer
	numberOfProblems int
}

func NewDiagnoseCmd() *cobra.Command {
	const (
		use   = "diagnose"
		short = "Diagnose the current configuration"
	)

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, _ []string) {
			diagnoser := diagnoser{
				ctx: cmd.Context(),
				out: cmd.OutOrStdout(),
			}
			diagnoser.diagnose()
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}

func (d *diagnoser) diagnose() {
	d.checkCommandEnv()
	if host, ok := d.checkHostEnv(); ok {
		d.checkConnection(host)
	}
	fmt.Fprintln(d.out, "Diagnostics are done.")
}

func (d *diagnoser) checkCommandEnv() {

	value, ok := os.LookupEnv("DOCKER_COMMAND")
	if !ok {
		d.reportBad("Environment variable DOCKER_COMMAND is not defined.")
		return
	}

	exePath, err := os.Executable()
	if err != nil && value != exePath {
		d.reportBad(`Environment variable DOCKER_COMMAND has a wrong value.
The actual value is '%s', but the full path of the current executable is '%s'
`,
			value, exePath)
		return
	}

	d.reportGood("Environment variable DOCKER_COMMAND has a valid value.")
}

func (d *diagnoser) checkHostEnv() (string, bool) {

	const (
		example = "tcp://192.168.0.100:2375"
	)

	host, ok := os.LookupEnv("DOCKER_HOST")
	if !ok {
		d.reportBad("Environment variable DOCKER_HOST is not defined.")
		return "", false
	}

	if err := docker.ValidateHost(host); err != nil {
		d.reportBad(
			`Environment variable DOCKER_HOST has an invalid value '%s': %s
A valid example is '%s' if the daemon is listening on the IP address and the port using TCP.`,
			host, err, example)
		return host, false
	}

	d.reportGood("Environment variable DOCKER_HOST has a valid value.")
	return host, true
}

func (d *diagnoser) checkConnection(host string) {

	fmt.Fprintf(d.out,
		"Verifying the connectivity to the Docker daemon at '%s' given by DOCKER_HOST...\n",
		host)

	c, err := docker.NewClientForHost(host)
	if err != nil {
		d.reportBad("Failed to create a client: %s", err)
		return
	}
	c.Close()

	_, err = c.Ping(d.ctx, client.PingOptions{})
	if err != nil {
		err = docker.WrapError(err)
		d.reportBad("The health check for the Docker daemon failed: %s", err)
		return
	}

	d.reportGood("Confirmed that the daemon is running in healthy state at the specified location.")
}

func (d *diagnoser) reportGood(format string, a ...any) {
	fmt.Fprintf(d.out, "✅"+format+"\n", a...)
}

func (d *diagnoser) reportBad(format string, a ...any) {
	fmt.Fprintf(d.out, "❌"+format+"\n", a...)
	d.numberOfProblems++
}
