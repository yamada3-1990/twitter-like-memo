※Qiitaに同じ記事があります!  
[GoアプリケーションのDocker imageをGitHub Actionsでビルド&プッシュする](https://qiita.com/yamada3/items/18e3cfebb5484f678470)  


CI/CAの勉強記録もかねてのメモ
## やりたいこと
タイトルの通り
Golangで開発したアプリケーションを、GitHub Actionsを使って自動的にDocker imageに変換し、コンテナレジストリにアップロードしたい

![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3953070/ecfc20b5-3348-43f8-89b9-3483b4efb730.png)
リポジトリトップページのPackagesのところ

## そもそもCI/CAとは
### CI（Continuous Integration：継続的インテグレーション）
>すべてのコード変更を共有ソースコードリポジトリのmainブランチに早い段階で頻繁に統合し、コミットやマージ時に各変更を自動的にテストし、自動的にビルドを開始するプラクティスのこと

* コードの自動ビルド
* 自動テスト
* 静的解析　など

任意のアクション(commit/push/PRが作られた時など)をトリガーとしてテストや検証を行うことで、早期に不具合を発見できたり、開発速度の向上につながる

### CD（Continuous Delivery/Continuous Deployment：継続的デリバリー/継続的デプロイメント）

>継続的デプロイとは、すべてのCI検証済みコードを本番環境に自動的にデプロイすることを意味し、一方で継続的デリバリーはこのコードをデプロイできることを意味します。

二つの概念があるらしい  
**継続的デリバリー**：デプロイできるようにソフトウェアを構築する  
**継続的デプロイメント**：コードの変更を自動的にデプロイまでする  

継続的デリバリー後に自動的にデプロイするようにすると継続的デプロイメントになる
継続的デリバリーが継続的デプロイメントの前提条件

## 本題
* リポジトリの「Settings」→「Actions」→「General」でGitHub Actionsが有効になっていることを確認

* ```.github/workflows/build.yml```を作成する

*余談 .ymlを作成したときに、GitHub Actionsという拡張機能おすすめ！とポップアップが出てきたため言われるがまま入れた VSC上でいろいろ確認できるらしい
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3953070/f3901757-84be-4a56-8efd-be827da1357b.png)


```build.yml
name: build
run-name: building ${{ github.ref_name }} now!
on: [push]
```

**name** : ワークフローの名前
**run-name** : ワークフローの実行タイトル
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3953070/a1f830d9-9294-4b6f-b2e8-164f9dd70dc4.png)

**on** : ワークフロー実行のトリガー 今回はpushをトリガーとした

```build.yml
env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}
```
**env** : 環境変数
**REGISTRY** : Docker imageをプッシュする先
**IMAGE_NAME** : Dockerイメージの名前



```build.yml
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
```
**runs-on** : GitHub Actionsの実行環境

```build.yml
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to the Container registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
```
**steps** : job内で実行されるステップリスト  
**name** : ステップの名前  
**uses** : 使用するGitHub Actionsのアクションを指定  

```actions/checkout```でリポジトリのコードをワークフローの実行環境にコピー  
↓  
```docker/setup-buildx-action```でDocker Buildxを設定  
↓  
```docker/login-action```でコンテナレジストリにログイン  

```build.yml
      - name: Extract metadata (tags, labels) for Docker
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=raw,value=latest
            type=ref,event=branch
            type=ref,event=tag
            type=ref,event=pr
            type=sha,prefix=,format=short
            
      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
**tags** : imageに付けるtagを指定  
**labels** : imageに付けるlabelを指定  
**cache-from** : ビルドキャッシュの読み込み元を指定 ここではGitHub Actionsのキャッシュから読み込むことでビルド時間の短縮に  
**cache-to** : ビルドキャッシュの保存先を指定 後続のビルドのcache-fromで再利用できるようにする  

```docker/metadata-action```でDockerイメージのタグとラベルを自動生成  
↓  
```docker/build-push-action```でDockerイメージをビルドし、生成されたタグとラベルを付与してコンテナレジストリにプッシュ。ビルドキャッシュをGitHub Actionsのキャッシュに保存し、後続のビルドで再利用できるようにする  

コード全体
```build.yml
name: build
run-name: building ${{ github.ref_name }} now!
on: [push]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to the Container registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata (tags, labels) for Docker
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=raw,value=latest
            type=ref,event=branch
            type=ref,event=tag
            type=ref,event=pr
            type=sha,prefix=,format=short
      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
- - -

変更をcommit&pushすると自動でビルド、プッシュされる
リポジトリトップページのpackageからアクセスすると
```
$ docker pull ghcr.io/github-id/~~~
```

みたいなコマンドがあるので実行する
適切にpullできたら
```
$ docker run -p 9000:9000 ghcr.io/github-id/~~~
```
でプルしてきたimageをrunして確認
※コード内で9000ポートを指定しているのでrunでも指定。指定していなかったら```-p```オプションは不要

## 参考
https://about.gitlab.com/ja-jp/topics/ci-cd/

https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#jobsjob_idenv

https://qiita.com/shun198/items/14cdba2d8e58ab96cf95#jobs