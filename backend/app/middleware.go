package app

import (
	"log/slog"
	"net/http"
	"strings"
)

// server.goのline:111で使われている

// MARK: - simpleCORSMiddleware()
// CORS(Cross-Origin Resource Sharing)ミドルウェア関数
func simpleCORSMiddleware(next http.Handler, origin string, methods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HTTPレスポンスのヘッダーにAccess-Control-Allow-Originを設定
		// どのオリジンからのリクエストを許可するか
		// origin; http://example.com とかの、プロトコル＋ホスト名＋ポート番号の組み合わせ
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// HTTPレスポンスのヘッダーにAccess-Control-Allow-Methodsを設定
		// 許可するメソッド(GETとかPOSTとか)を指定
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		// HTTPレスポンスのヘッダーにAccess-Control-Allow-Headersを設定
		// 実際の要求で使用できるカスタムヘッダーをブラウザに伝える
		// ここではすべてのヘッダーを許可するためにワイルドカード(*)を使用
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// プリフライトリクエスト: 異なるオリジン間のリソース共有(CORS)を行う際に、ブラウザが自動的に送信する事前確認用のHTTPリクエスト
		// HTTPメソッドがOPTIONSだったら=サーバーがリクエストを受け入れられる場合
		if r.Method == "OPTIONS" {
			// 成功のレスポンスを返す
			w.WriteHeader(http.StatusOK)
			return
		}

		// HTTPメソッドがOPTIONSじゃなかったら
		// リクエストを次のハンドラーに渡す
		next.ServeHTTP(w, r)
	})
}

// MARK: - simpleLoggerMiddleware()
// HTTPリクエストがサーバーに届いた際に、そのリクエストに関する基本的な情報をログに出力するための関数
// next http.Handler: ロギング処理が終わった後にリクエストを渡す先のハンドラー
func simpleLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// .UserAgent(): returns the client's User-Agent
		// ユーザーエージェント: 私(クライアント)は〇〇という種類のソフトウェアで、△△という環境で動いていますよ～という情報
		slog.Info("request received", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr, "user_agent", r.UserAgent())
		next.ServeHTTP(w, r)
	})
}
