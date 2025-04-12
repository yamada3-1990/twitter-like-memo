package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var errMemoNotFound = errors.New("memo not found")

// Memo構造体
type Memo struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Body  string `db:"body" json:"body"`
	Tags  string `db:"tags" json:"tags"`
}

// ItemRepositoryインターフェース
// どのようなメソッドを持つかを定義
type ItemRepository interface {
	Insert(ctx context.Context, memo *Memo) error
	GetAllMemos(ctx context.Context) ([]Memo, error)
	Delete(ctx context.Context, memo *Memo) error
	SearchByKeyword(ctx context.Context, keyword string) ([]Memo, error)
	SearchByTags(ctx context.Context, tags []string) ([]Memo, error)
}

// itemRepository構造体
// ItemRepositoryインターフェースを実際に実装する型
type itemRepository struct {
	// sql.DB型のポインタを保持
	db *sql.DB
}

// MARK: - NewItemRepository()
// 新しいitemRepositoryインスタンスを作成
// *sql.DB(データベース接続)を受け取る
func NewItemRepository(db *sql.DB) (ItemRepository, error) {
	// スキーマを読み込んでテーブルが無かったら作成
	q, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(q))
	if err != nil {
		slog.Error("failed to create memos table", "error", err)
		return nil, err
	}

	// 作成したitemRepositoryのインスタンスを返す
	return &itemRepository{db: db}, nil
}

// MARK: - Insert()
// 新しいメモをデータベースに保存
func (i *itemRepository) Insert(ctx context.Context, memo *Memo) error {
	// トランザクション(複数の操作をまとめて実行)を開始
	// 全部成功したら変更を適用(コミット)、失敗したら破棄(ロールバック)
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// titleとbodyをインサート
	query := `INSERT INTO memos (title, body) VALUES (?, ?)`
	// ExecContext: 行を返さずにクエリを実行
	res, err := tx.ExecContext(ctx, query, memo.Title, memo.Body)
	if err != nil {
		return err
	}

	// 今インサートしたメモのidを取得
	memoID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// タグ処理
	// カンマ区切りで分割
	tagNames := strings.Split(memo.Tags, ",")

	for _, tagName := range tagNames {
		//前後の空白を削除
		tagName = strings.TrimSpace(tagName)
		// tagsテーブルに同じタグがあるかを確認
		var tagID int64
		// QueryRow: 最大1行を返すと予想されるクエリを実行
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// タグが存在しない場合はインサート
			result, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
			if err != nil {
				return err
			}
			tagID, err = result.LastInsertId()
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// メモとタグの関連付け
		_, err = tx.Exec("INSERT INTO memo_tags (memo_id, tag_id) VALUES (?, ?)", memoID, tagID)
		if err != nil {
			return err
		}
	}

	// 全ての処理が成功したらトランザクションをコミット
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// MARK: - GetAllMemos()
// データベースから全てのメモを取得
func (i *itemRepository) GetAllMemos(ctx context.Context) ([]Memo, error) {
	query := `
				SELECT
					memos.*,
					COALESCE(GROUP_CONCAT(tags.name), '') AS tags
				FROM
					memos
				LEFT JOIN
					memo_tags ON memos.id = memo_tags.memo_id
				LEFT JOIN
					tags ON memo_tags.tag_id = tags.id
				GROUP BY
					memos.id;
			`
	// QueryContext: 行を返すクエリを実行
	rows, err := i.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []Memo
	// .Next(): 次の行があるかを確認(bool)
	for rows.Next() {
		var memo Memo
		// .Scan(): 一致した行の列を引数が指す値にコピー
		err := rows.Scan(&memo.ID, &memo.Title, &memo.Body, &memo.Tags)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

// MARK: - Delete()
// 指定されたメモをデータベースから削除
func (i *itemRepository) Delete(ctx context.Context, memo *Memo) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// まずメモとタグの関連付けを削除
	_, err = tx.ExecContext(ctx, "DELETE FROM memo_tags WHERE memo_id = ?", memo.ID)
	if err != nil {
		return err
	}

	// メモを削除
	query := `DELETE FROM memos WHERE id = ?`
	result, err := tx.ExecContext(ctx, query, memo.ID)
	if err != nil {
		return err
	}

	// RowsAffected: update, insert, deleteによって影響を受けた行数を返す
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("memo not exist")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// MARK: - SearchByKeyword()
// 指定されたkeywordがtitleまたはbodyに含まれるメモを検索
func (i *itemRepository) SearchByKeyword(ctx context.Context, keyword string) ([]Memo, error) {
	query := `
				SELECT
					memos.id,
					memos.title,
					memos.body,
					GROUP_CONCAT(tags.name) AS tags
				FROM
					memos
				LEFT JOIN
					memo_tags ON memos.id = memo_tags.memo_id
				LEFT JOIN
					tags ON memo_tags.tag_id = tags.id
				WHERE
					memos.title LIKE '%' || ? || '%' OR memos.body LIKE '%' || ? || '%'
				GROUP BY
					memos.id;
			`
	// Query: 行を返すクエリを実行
	rows, err := i.db.Query(query, keyword, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []Memo
	for rows.Next() {
		var memo Memo
		err := rows.Scan(&memo.ID, &memo.Title, &memo.Body, &memo.Tags)
		if err != nil {
			// var ErrNoRows = errors.New("sql: no rows in result set")
			if err == sql.ErrNoRows {
				return []Memo{}, errMemoNotFound
			} else {
				return nil, err
			}
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

// MARK: - SearchByTags()
// 指定されたtagを持つメモを検索
// tagsはtag1,tag2,...の形式だけど、server.goで呼び出される前に
// parseSearchByTagsRequest()で"tag1", "tag2",...の形式に変換される
func (i *itemRepository) SearchByTags(ctx context.Context, tags []string) ([]Memo, error) {
	// tagが無かったら空のスライスを返す
	if len(tags) == 0 {
		return []Memo{}, nil
	}

	// tagsスライスと同じ長さの文字列スライスを作成
	subqueries := make([]string, len(tags))

	// tagsスライスと同じ長さのインターフェーススライスを作成
	// SQLのプレースホルダの?に渡す引数を格納
	args := make([]interface{}, len(tags))
	for i, tag := range tags {
		// orじゃなくてand検索にしている
		subqueries[i] = `EXISTS (
			SELECT 1 FROM memo_tags mt
			JOIN tags t ON mt.tag_id = t.id
			WHERE mt.memo_id = memos.id AND t.name = ?
		)`
		args[i] = tag
	}

	// 全てのサブクエリをANDで結合
	query := `
		SELECT DISTINCT
			memos.id,
			memos.title,
			memos.body,
			GROUP_CONCAT(DISTINCT tags.name) AS tags
		FROM
			memos
		LEFT JOIN
			memo_tags ON memos.id = memo_tags.memo_id
		LEFT JOIN
			tags ON memo_tags.tag_id = tags.id
		WHERE
			` + strings.Join(subqueries, " AND ") + `
		GROUP BY
			memos.id, memos.title, memos.body;
	`

	// // デバッグ用
	// slog.Info("executing query", "query", query, "args", args)

	rows, err := i.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("query execution failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var memos []Memo
	for rows.Next() {
		var memo Memo
		err := rows.Scan(&memo.ID, &memo.Title, &memo.Body, &memo.Tags)
		if err != nil {
			slog.Error("row scanning failed", "error", err)
			return nil, err
		}
		memos = append(memos, memo)
	}
	// .Err():  Row.Scan.Errを呼び出さないでもクエリエラーをチェックできる
	if err = rows.Err(); err != nil {
		slog.Error("rows iteration failed", "error", err)
		return nil, err
	}
	// デバッグ用
	// slog.Info("query completed", "result_count", len(memos))
	return memos, nil
}
