# バックエンドのビルドステージ
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY backend ./backend
COPY db/schema.sql ./db/schema.sql

RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 go build -o main ./backend/cmd

# フロントエンドのビルドステージ
FROM node:22.12.0-alpine AS frontend-builder

WORKDIR /app

COPY frontend/twitter-like-memo/package*.json ./
RUN npm install

COPY frontend/twitter-like-memo/ ./
ENV VITE_API_URL=/api
RUN npm run build

# 実行用の最終ステージ
FROM alpine:3.19

WORKDIR /app

# バックエンドの実行に必要なライブラリをインストール
RUN apk add --no-cache libc6-compat

# nginx のインストールと設定
RUN apk add --no-cache nginx
COPY frontend/twitter-like-memo/nginx.conf /etc/nginx/http.d/default.conf

# バックエンドのバイナリをコピー
COPY --from=backend-builder /app/main ./
COPY --from=backend-builder /app/db/schema.sql ./db/schema.sql

# フロントエンドのビルド成果物をコピー
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# ポートの公開
EXPOSE 80 9000

# 起動スクリプト
COPY start.sh ./
RUN chmod +x start.sh

CMD ["./start.sh"]