package app

import (
	"log/slog"
	"os"
)

type Server struct {
	// ポート番号
	Port string
	// 画像を保存するディレクトリへのパス
	ImageDirPath string
}

func (s Server) Run() int {
	opts := slog.HandlerOptions{
		level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	slog.SetDefault(logger)
}
