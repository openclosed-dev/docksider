# Docksider

A Docker-style CLI for pulling and pushing images directly to container registries, not via Docker daemon.

## Configuration

Set up the environment variables in Windows.

| name | value |
| - | - |
| DOCKER_COMMAND | the full path to this executable including filename |
| DOCKER_HOST | the URL of the Docker daemon specified in the format `tcp://<address>:<port>` |

Note that the Docker daemon must be up and running on the specified address and port.

See [Configure remote access for Docker daemon](https://docs.docker.com/engine/daemon/remote-access/)

## Using with Azure Container Registry

1. Log in to your Azure container registry using Azure CLI.

    ```shell
    az login
    az acr login -n <registry>
    ```

2. Upload a container image to the container registry, which is retrieved from the Docker daemon.

    ```shell
    docksider push <registry>/<image>:<tag>
    ```
