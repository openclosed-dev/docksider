# Docksider

[![Build](https://github.com/openclosed-dev/docksider/actions/workflows/build.yml/badge.svg)](https://github.com/openclosed-dev/docksider/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/release/openclosed-dev/docksider/all.svg)](https://github.com/openclosed-dev/docksider/releases)


Dockerデーモン経由ではなく、コンテナレジストリにイメージを直接プルおよびプッシュするためのDockerスタイルのCLI。

## 特徴

* Docker デーモン経由ではなく、コンテナイメージをコンテナレジストリに直接プルおよびプッシュします。
* Docker CLI コマンドのサブセットを提供します。
* Azure CLI と連携してコンテナレジストリに認証します。

## 使い方

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

## インストール

Windows用の実行ファイルは、[Releases](https://github.com/openclosed-dev/docksider/releases) ページからダウンロードでき、`PATH` 環境変数に含まれるフォルダに保存する必要があります。

あるいは、Go言語のSDKがインストールされている場合は、次のコマンドを実行するだけで、このプログラムをローカル環境にインストールできます。

```shell
go install github.com/openclosed-dev/docksider/cmd/docksider@latest
```

## 設定

Windowsで次の環境変数を設定します。

### `DOCKER_HOST`

Docker デーモンの URL。形式は`tcp://<address>:<port>`。

指定されたアドレスとポートで Dockerデーモンが起動して実行されている必要があることに注意してください。

[Configure remote access for Docker daemon](https://docs.docker.com/engine/daemon/remote-access/)を参照してください。

### `DOCKER_COMMAND`

この実行可能ファイルへのフルパス。ファイル名を含みます。
この変数は、Azure CLIがこのプログラムを正しく検出するために必要です。

## Azure Container Registryと使う

1. Azure CLIを使用して Azureコンテナレジストリにログインします。

    ```shell
    az login
    az acr login -n <registry>
    ```

2. Dockerデーモンから取得したコンテナ イメージをコンテナレジストリにアップロードします。

    ```shell
    docksider push <registry>.azurecr.io/<image>:<tag>
    ```
