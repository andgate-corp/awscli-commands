# awscli-commands
Interactive command in slack for aws-cli


# 開発環境
- go 1.x
    - REQUIRE dep
- Node.js
    - REQUIRE Serverless framework
- Docker
    - 動作確認済みの環境: Docker version 18.06.0-ce, build 0ffa825

# 事前設定
- Slack Appの設定
    - Slash Commands の有効化
    - Interactive Components の有効化
    - Verification Tokenを取得しておくこと

- ディレクトリ内にcustom.ymlを作成する
```custom.yml
verificationToken: <Slack AppのVerificationToken>
```


# デプロイ

- Dockerの起動

dockerデーモンを起動して下さい

- golang用コンテナの作成

golangをbuildするためのコンテナ(awscli-command:latest)を生成します.

```
make build-docker-golang
```

- 依存パッケージの導入

```
make run-dep-ensure
```

- ビルド

```
make run-build
```

- デプロイ

softinstigate/serverless へ `$HOME/.aws` と このDirectoryをVOLUME共有して `sls deploy` を実行します。

```
make sls-deploy aws_profile=<deployターゲットのAWS PROFILE名> region=<AWS REGION>
```

# ビルド後設定
- Slach Commands のコマンドを追加する
    - Request URLに`https://<endpoint>/<stage>/awscli` のURLを設定する
- Interactive Componentsを設定する
    - Interactivity のRequestURLに`https://<endpoint>/<stage>/awscli/interaction` のURLを設定する
    - Actionsは不要

# 構造
```
awscli-commands  
    |- commands             ->  awscli command codes
    |- handler              ->  lambda handlers
    |   |- igniter          ->  Slash command handler
    |   |- interaction      ->  Interactive Components handler
    |
    |- serverless.yml       ->  Serverless config file
    |- custom.yml           ->  Credential config file
    |- Gopkg.lock
    |- Gopkg.toml
    |- Makefile             -> build commands
    |- README.md
```