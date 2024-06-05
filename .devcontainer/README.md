# 開発環境について

## devcontainer により開発環境を開く方法

この環境では、vscode の devcontainer という機能を利用し、VS Code でコンテナ内に入って作業することができる。
devcontainer の詳細については、[devcontainer を使用した開発についての公式ドキュメント](https://code.visualstudio.com/docs/devcontainers/containers)を参照すること。

### 前提
- gitによるリポジトリの同期
- Docker のインストール (Windows の場合は WSL2 上にインストール)
- Docker/Rancher Desktop for Windows/Mac のインストール
- VS Code のインストール
- VS Code の拡張機能 - Dev Containers のインストール

コンテナの起動が成功することを確認しているバージョン

- Dev Containers - version 0.338.1
- VS Code - version 1.86.0
  ※ 筆者の環境では、VS Code と Dev Containers の環境が最新でない場合に、コンテナの起動が失敗する場合があった。

### 操作
**プラグインの開発環境コンテナに入る手順**
- VS Code で .devcontainer を含むディレクトリを開き、左下の「><」をクリックし、以下のメニューを表示する。
  ![リモートコンテナに接続するメニューを開く様子のキャプチャ](./images/vscode-capture-open-remote-menu.png)
- 次に、[コンテナーで再度開く] を押下する。
  ![コンテナで再度開く](./images/image.png)
- すると、プラグインの開発環境である[gf-dev container]と、grafana のコンテナである[grafana container]が起動し、開発環境[gf-dev container]の内部に入る。

## grafana plugin 開発での操作

- [gf-dev container]を起動し、新たなターミナルを開くと、コンテナ内の /workspace ディレクトリが開く。
- `cd sios-fiap-datasource`コマンドで、grafana plguin 開発用作業ディレクトリへ移動する。

### Grafanaへのアクセス

ビルドしたバイナリを配置し、動作させるGrafanaへのアクセスは、`localhost:3000`で行う。

### 開発時の操作

フロントエンドとバックエンドのビルドを行うためのコマンドは以下の通り。
※ビルドを Grafana に反映させるためには、grafana の service を再起動する必要がある。

```bash
# 依存関係のインストール
npm install
# プラグインのフロントエンドをビルド(開発モード)
npm run dev
# プラグインのフロントエンドをビルド
npm run build
# 新しいターミナルで、バックエンドをビルドする
mage -v build:linux
```

### Grafana の再起動方法

docker-compose.yml の存在するディレクトリで、外部のターミナルから`docker compose restart grafana`を実行。  
ビルド後に、この作業を行うことで、Grafanaに変更が反映される。

## コンテナの構成についての説明

`docker-compose.yml`で、grafana と dev という 2 つの service を定義している。

- grafana
  - grafana を起動する
    - 環境変数で Grafana の設定を変更可能
    - Grafana イメージをベースに作成した Dockerfile を使って環境を作成している
- dev
  - grafana backend datasource plugin を開発するための環境を作成する
    - Ubuntu イメージをベースとして必要な環境をセットアップ
      - セットアップ内容
        - nodejs 20.9.0
        - go 1.21.7
        - mage 1.15.0