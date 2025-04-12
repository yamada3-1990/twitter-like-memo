package main

import (
	"os"
	// appパッケージ(apiが実装されているパッケージ)をインポート
	"twitter-like-memo/backend/app"
)

const (
	// backendのポート
	port         = "9000"
	imageDirPath = "images"
)

func main() {
	// Exit: プログラムを終了させる
	os.Exit(app.Server{
		// server.goのServer構造体のインスタンスを作成
		Port:         port,
		ImageDirPath: imageDirPath,
	}.Run())
}
