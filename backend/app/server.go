package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	// ItemRepositoryインターフェース
	itemRepo ItemRepository
}

// AddMemoRequest 構造体
type AddMemoRequest struct {
	ID    int    `form:"id"`
	Title string `form:"title"`
	Body  string `form:"body"`
	Tags  string `form:"tags"`
}

// メモ追加時に返されるメッセージ
type AddMemoResponse struct {
	Message string `json:"message"`
}

// keyword検索用の構造体
type SearchByKeywordRequest struct {
	keyword string
}

// tag検索用の構造体 複数のtagを格納するためにスライスを利用
type SearchByTagsRequest struct {
	tags []string
}

// MARK: - Run()
/*
ログ設定
CORS設定
データベースへの接続
ハンドラ設定
ルーティング
サーバの起動
*/
func (s Server) Run() int {
	// ログ設定
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	slog.SetDefault(logger)

	// CORSの設定
	// フロントエンドのURL
	// .LookupEnv(): LookupEnv(key) keyで指定された環境変数を取得 存在する場合はその値とtrueを返す
	frontURL, found := os.LookupEnv("FRONT_URL")
	if !found {
		frontURL = "http://localhost:5173"
	}

	// .Open(): dbドライバとデータソース名を指定して開く
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
	// Handlers構造体の新しいインスタンスを作成する
	h := &Handlers{imgDirPath: s.ImageDirPath, itemRepo: itemRepo}

	// ルーティング
	// .NewServeMux(): 新しいServeMuxを割り当てて返す
	// ServeMux: HTTPリクエストのマルチプレクサ(multiplexer, 複数の入力信号を1つの出力信号に合成する論理回路)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /memos", h.AddMemo)
	mux.HandleFunc("GET /memos", h.GetMemos)
	mux.HandleFunc("DELETE /memos/{id}", h.DeleteMemo)
	mux.HandleFunc("GET /search/keyword", h.SearchByKeyword)
	mux.HandleFunc("GET /search/tags", h.SearchByTags)

	// サーバーの起動
	slog.Info("http server started on", "port", s.Port)
	// .ListenAndServe(): 与えられたアドレスとハンドラでHTTPサーバを開始
	// simpleCORSMiddleware(): middleware.goで定義されている
	// simpleLoggerMiddleware(): middleware.goで定義されている
	err = http.ListenAndServe(":"+s.Port, simpleCORSMiddleware(simpleLoggerMiddleware(mux), frontURL, []string{"GET", "HEAD", "POST", "DELETE", "OPTIONS"}))
	if err != nil {
		slog.Error("failed to start server: ", "error", err)
		return 1
	}

	return 0
}

// MARK: - parseAddMemoRequest()
// Memoの追加リクエストをパースする
// HTTPリクエスト(*http.Request)からメモの情報を取り出してAddMemoRequest構造体に格納
func parseAddMemoRequest(r *http.Request) (*AddMemoRequest, error) {
	// マルチパートフォームデータを解析
	// .ParseMultipartForm(): リクエストボディをmultipart/form-dataとして解析する
	err := r.ParseMultipartForm(10 << 20) // 10MBの制限
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err)
		return nil, err
	}

	// フォームデータの内容をログ出力
	slog.Info(
		"received form data",
		"form", r.Form,
		"postForm", r.PostForm,
		"multipartForm", r.MultipartForm,
	)

	// HTTPリクエストのフォームデータからtagsという名前のフィールドの値を取得
	tags := r.Form["tags"]
	var tagList []string
	for _, tag := range tags {
		// カンマ区切りのタグを分割
		// .Split(s, sep): sをsepで分割された部分文字列にスライス
		tagList = append(tagList, strings.Split(tag, ",")...)
	}

	req := &AddMemoRequest{
		// .FormValue(): クエリの値を返す
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
		Tags:  strings.Join(tagList, ","),
	}

	// パースされたリクエストの内容をログ出力
	slog.Info(
		"parsed request",
		"title", req.Title,
		"body", req.Body,
		"tags", req.Tags,
	)

	// バリデーション
	if req.Title == "" {
		req.Title = "無題"
	}

	if req.Body == "" {
		return nil, errors.New("body is required")
	}

	return req, nil
}

// MARK: - AddMemo()
// POST /memos でメモを追加
func (s *Handlers) AddMemo(w http.ResponseWriter, r *http.Request) {
	// リクエストのコンテクストを返す
	ctx := r.Context()

	// リクエストをパースする
	req, err := parseAddMemoRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memo := &Memo{
		Title: req.Title,
		Body:  req.Body,
		Tags:  req.Tags,
	}
	message := fmt.Sprintf("memo received: %s", memo.Title)
	slog.Info(message)

	err = s.itemRepo.Insert(ctx, memo)
	if err != nil {
		slog.Error("failed to store memo: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Memo追加が成功したメッセージを追加
	resp := AddMemoResponse{Message: message}
	// .NewEncoder(): wに書き込む新しいエンコーダを返す
	// .Encode(): JSON形式にエンコーディングする
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		// StatusInternalServerError: 500 ウェブサイトをホストしているサーバー側で何らかの不具合が発生
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// MARK: - GetMemos()
// GET /memos でメモを取得
func (s *Handlers) GetMemos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	memos, err := s.itemRepo.GetAllMemos(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		// 複数のMemoオブジェクトを格納するため
		Memos []Memo `json:"memos"`
	}{
		Memos: memos,
	}

	// w.Header()でレスポンスのHTTPヘッダーを取得する
	// 取得したHTTPヘッダーに対して、Content-Type ヘッダーを設定、値は"application/json"
	// これから送信するレスポンスボディはJSON形式だよーと通知する
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// MARK: - DeleteMemo()
// DELETE /memos/{id} でメモを削除
func (s *Handlers) DeleteMemo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URLからIDを取得
	// .PathValue(): リクエストにマッチしたServeMuxパターン中の指定されたパスワイルドカードの値を返す
	// "id"という名前のパスワイルドカードの値を取得している
	// Path Wildcard: URLの特定の部分を変数として扱うための記法 /memos/{id}の{}の部分
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// IDを整数に変換
	memoID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Memo構造体の新しいインスタンスを作成
	memo := &Memo{
		// IDフィールドのみを初期化
		ID: memoID,
	}
	message := fmt.Sprintf("deleted memo: %d", memo.ID)
	slog.Info(message)

	err = s.itemRepo.Delete(ctx, memo)
	if err != nil {
		if errors.Is(err, errors.New("memo not exist")) {
			slog.Error("memo not exist", "error", err)
			// StatusNotFound: 404 ページが存在しない
			http.Error(w, "memo not exist", http.StatusNotFound)
			return
		}
		slog.Error("failed to delete memo: ", "error", err)
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

// MARK: - parseSearchByKeywordRequest()
// Memoのキーワード検索リクエストをパースする
func parseSearchByKeywordRequest(r *http.Request) (*SearchByKeywordRequest, error) {
	// keyword用の構造体に、クエリパラメータからkeywordを取得して格納
	// 'http://127.0.0.1:9000/search/keyword?keyword=おはよう'の ? で始まる部分がクエリパラメータ
	req := &SearchByKeywordRequest{
		keyword: r.URL.Query().Get("keyword"),
	}

	// バリデーション
	if req.keyword == "" {
		return nil, errors.New("keyword is required")
	}
	return req, nil
}

// MARK: - SearchByKeyword()
// GET /search/keyword でメモを検索
func (s *Handlers) SearchByKeyword(w http.ResponseWriter, r *http.Request) {
	req, err := parseSearchByKeywordRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// SearchByKeyword関数にkeywordを渡す
	memos, err := s.itemRepo.SearchByKeyword(r.Context(), req.keyword)

	if err != nil {
		if errors.Is(err, errMemoNotFound) {
			slog.Warn("memo not exist: ", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if memos == nil {
		memos = []Memo{}
	}

	// .Marchal(): JSONのエンコーディングを返す
	jsonData, err := json.Marshal(memos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// JSONにエンコーディングした結果をレスポンスボディに書き込む
	w.Write(jsonData)
}

// MARK: - parseSearchByTagsRequest()
// Memoのタグ検索リクエストをパースする
func parseSearchByTagsRequest(r *http.Request) (*SearchByTagsRequest, error) {
	// クエリパラメータから直接取得
	tags := r.URL.Query()["tags"]
	if len(tags) == 0 {
		return nil, errors.New("tags is required")
	}

	var tagList []string
	for _, tag := range tags {
		// カンマ区切りのタグを分割
		tagList = append(tagList, strings.Split(tag, ",")...) 
	}

	req := &SearchByTagsRequest{
		tags: tagList,
	}

	return req, nil
}

// MARK: - SearchByTag()
// GET /search/tags でメモを検索
func (s *Handlers) SearchByTags(w http.ResponseWriter, r *http.Request) {
	req, err := parseSearchByTagsRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memos, err := s.itemRepo.SearchByTags(r.Context(), req.tags)

	if err != nil {
		if errors.Is(err, errMemoNotFound) {
			slog.Warn("memo not exist: ", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if memos == nil {
		memos = []Memo{}
	}

	jsonData, err := json.Marshal(memos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
