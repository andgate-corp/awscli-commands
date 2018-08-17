# awscli-commands
Interactive command in slack for aws-cli


# 開発環境
- go 1.x
    - REQUIRE dep
- Node.js
    - REQUIRE Serverless framework



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
- ビルド
```
make build
```

- デプロイ
```
sls deploy
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