FROM golang:1.26.0-trixie

ARG USERNAME
ARG USER_UID
ARG USER_GID=$USER_UID
ARG DOCKER_GID

RUN groupadd -g ${DOCKER_GID} docker

RUN groupadd -g ${USER_GID} ${USERNAME} \
  && useradd -u ${USER_UID} -g ${USER_GID} -m -s /bin/bash ${USERNAME} \
  && usermod -aG docker ${USERNAME}

RUN apt-get update \
  && export DEBIAN_FRONTEND=noninteractive \
  && apt-get install -y --no-install-recommends git openssh-client

# Install Azure CLI
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

RUN mkdir /workspaces && chown ${USERNAME}:${USERNAME} /workspaces

USER ${USERNAME}

RUN mkdir ~/.azure

RUN go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install golang.org/x/tools/gopls@latest
