package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	// ポート番号
	Port string
	// 画像を保存するディレクトリへのパス
	ImageDirPath string
}

type Handlers struct {
	// 画像を保存するディレクトリへのパス
	imgDirPath string
	itemRepo   ItemRepository
}

type AddMemoRequest struct {
	ID    int    `form:"id"`
	Title string `form:"title"`
	Body  string `form:"body"`
}

type AddMemoResponse struct {
	Message string `json:"message"`
}

func (s Server) Run() int {
	// ログ設定
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	slog.SetDefault(logger)

	// CORSの設定
	frontURL, found := os.LookupEnv("FRONT_URL")
	if !found {
		frontURL = "http://localhost:3000"
	}

	db, err := sql.Open("sqlite3", "db/memo.sqlite3")
	if err != nil {
		slog.Error("failed to open database: ", "error", err)
		return 1
	}
	defer db.Close()

	// handlerの設定
	itemRepo, err := NewItemRepository(db)
	if err != nil {
		slog.Error("failed to create item repository: ", "error", err)
		return 1
	}
	h := &Handlers{imgDirPath: s.ImageDirPath, itemRepo: itemRepo}

	// ルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("POST /memos", h.AddMemo)
	mux.HandleFunc("GET /memos", h.GetMemos)

	// サーバーの起動
	slog.Info("http server started on", "port", s.Port)
	err = http.ListenAndServe(":"+s.Port, simpleCORSMiddleware(simpleLoggerMiddleware(mux), frontURL, []string{"GET", "HEAD", "POST", "OPTIONS"}))
	if err != nil {
		slog.Error("failed to start server: ", "error", err)
		return 1
	}

	return 0
}

// Memoの追加リクエストをパースする
func parseAddMemoRequest(r *http.Request) (*AddMemoRequest, error) {
	var req = &AddMemoRequest{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}

	// バリデーション
	if req.Title == "" {
		req.Title = "無題"
	}

	if req.Body == "" {
		return nil, errors.New("body is required")
	}

	return req, nil
}

// POST /memos でメモを追加
func (s *Handlers) AddMemo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := parseAddMemoRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memo := &Memo{
		Title: req.Title,
		Body:  req.Body,
	}
	message := fmt.Sprintf("memo received: %s", memo.Title)
	slog.Info(message)

	err = s.itemRepo.Insert(ctx, memo)
	if err != nil {
		slog.Error("failed to store memo: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := AddMemoResponse{Message: message}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Handlers) GetMemos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	memos, err := s.itemRepo.GetAllMemos(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Memos []Memo `json:"memos"`
	}{
		Memos: memos,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
