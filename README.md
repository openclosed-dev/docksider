# Docksider

[![Build](https://github.com/openclosed-dev/docksider/actions/workflows/build.yml/badge.svg)](https://github.com/openclosed-dev/docksider/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/release/openclosed-dev/docksider/all.svg)](https://github.com/openclosed-dev/docksider/releases)

A Docker-style CLI for pulling and pushing images directly to container registries, not via Docker daemon.

[日本語](./README_ja.md)

## Features

* Pulls and pushes container images directly to a container registry instead of via the Docker daemon.
* Provides a subset of the Docker CLI commands.
* Works with Azure CLI to authenticate to container registries.

## Usage

```
A Docker-style CLI for pulling and pushing images directly to container registries

Usage:
  docksider [OPTIONS] COMMAND [ARG...]
  docksider [command]

Available Commands:
  help        Help about any command
  image       Manage images
  images      List images
  login       Authenticate to a registry
  pull        Download an image from a registry
  push        Upload an image to a registry

Flags:
  -h, --help   help for docksider

Use "docksider [command] --help" for more information about a command.
```

## Installation

The executable for Windows platforms can be downloaded from the [Releases](https://github.com/openclosed-dev/docksider/releases) page and should be saved in a folder included in your `PATH` environment variable.

Alternatively, if you have the Golang SDK installed, you can simply run the following command to install this program in your local environment.

```shell
go install github.com/openclosed-dev/docksider/cmd/docksider@latest
```

## Configuration

Set the following environment variables on Windows platforms.

### `DOCKER_HOST`

The URL of the Docker daemon in the format `tcp://<address>:<port>`.

Note that the Docker daemon must be up and running on the specified address and port.

See [Configure remote access for Docker daemon](https://docs.docker.com/engine/daemon/remote-access/)

If your daemon is running on a WSL distro, try running the following script from this repository on your WSL distro:

```shell
sudo bash configure-docker-daemon.sh
```

### `DOCKER_COMMAND`

The full path to this executable, including the filename.
This variable is required for Azure CLI to properly find this program.

## Using with Azure Container Registry

1. Log in to your Azure container registry using Azure CLI.

    ```shell
    az login
    az acr login -n <registry>
    ```

2. Upload the container image obtained from the Docker daemon to the container registry.

    ```shell
    docksider push <registry>.azurecr.io/<image>:<tag>
    ```
