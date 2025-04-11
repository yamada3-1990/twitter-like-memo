#!/bin/sh

# バックエンドを起動（バックグラウンドで）
./main &

# nginxをフォアグラウンドで起動
nginx -g 'daemon off;'
