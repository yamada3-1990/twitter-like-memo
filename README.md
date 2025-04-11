# twitter-like-memo

<p style="display: inline">
  <!-- フロントエンド -->
  <img src="https://img.shields.io/badge/-Node.js-000000.svg?logo=node.js">
  <img src="https://img.shields.io/badge/-React-20232A?logo=react">
  <img src="https://img.shields.io/badge/-Vite-646CFF?logo=vite">
  <!-- バックエンド -->
  <img src="https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white">
  <!-- ミドルウェア -->
  <img src="https://img.shields.io/badge/-Nginx-269539.svg?logo=nginx">
  <img src="https://img.shields.io/badge/-SQLite-003B57?logo=sqlite">
  <!-- インフラ -->
  <img src="https://img.shields.io/badge/-Docker-2496ED?logo=docker&logoColor=white">
  <img src="https://img.shields.io/badge/-GitHub_Actions-2088FF?logo=github-actions&logoColor=white">
</p>

TwitterライクなUIのメモアプリです。復習も兼ねて作成。

## 機能

- メモの作成
    - **title**: NULLの場合でもメモを追加できます。その場合は「無題」となります。
    - **body**: NULLの場合はメモを追加できません。
    - **tag**: ```tag1,tag2,...```のように追加してください。複数追加可能。

 - メモの検索
    - **キーワード検索**: titleとbodyにキーワードが含まれるメモを取得します。
    - **タグ検索**: ```tag1,tag2,...```のように指定してください。複数指定可能。
    - **結果をクリア**をクリックすることで、検索結果をクリアすることができます。

 - メモの削除
    - 各メモ右上のごみ箱アイコンをクリックしてください。削除時の確認メッセージは出ません。


## 使用したもの

### フロントエンド
- React 19
- React Router v7
- Vite

### バックエンド
- Go
- SQLite

### インフラ
- Docker
- Nginx 
- GitHub Actions

## 実行方法

### Dockerを使用して実行

1. imageのプル(Packagesのコマンドを参照してください)
```bash
docker pull ghcr.io/yamada3-1990/twitter-like-memo:☆☆☆
```

2. コンテナの起動
```bash
docker run -p 80:80 -p 9000:9000 twitter-like-memo
```

3. アクセス
   - http://localhost
   <!-- - バックエンドAPI: http://localhost:9000 -->

### リポジトリをcloneして確認もできます

Docker Composeを利用
```bash
docker-compose up
docker-compose down
```
