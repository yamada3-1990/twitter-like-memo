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

type Memo struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Body  string `db:"body" json:"body"`
	Tags  string `db:"tags" json:"tags"`
}

type ItemRepository interface {
	Insert(ctx context.Context, memo *Memo) error
	GetAllMemos(ctx context.Context) ([]Memo, error)
	Delete(ctx context.Context, memo *Memo) error
	SearchByKeyword(ctx context.Context, keyword string) ([]Memo, error)
	SearchByTags(ctx context.Context, tags []string) ([]Memo, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	db *sql.DB
}

// 新しいitemRepositoryを作成
func NewItemRepository(db *sql.DB) (ItemRepository, error) {
	// テーブルが無かったら作成
	q, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(q))
	if err != nil {
		slog.Error("failed to create memos table", "error", err)
		return nil, err
	}

	return &itemRepository{db: db}, nil
}

// MARK: - Insert()
func (i *itemRepository) Insert(ctx context.Context, memo *Memo) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// titleとbodyをインサート
	query := `INSERT INTO memos (title, body) VALUES (?, ?)`
	res, err := tx.ExecContext(ctx, query, memo.Title, memo.Body)
	if err != nil {
		return err
	}

	// 今インサートしたメモのidを取得
	memoID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// カンマ区切りで分割
	tagNames := strings.Split(memo.Tags, ",")

	for _, tagName := range tagNames {
		//前後の空白を削除
		tagName = strings.TrimSpace(tagName)
		// タグの存在を確認
		var tagID int64
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// タグが存在しない場合は挿入
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// MARK: - GetAllMemos()
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
	rows, err := i.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []Memo
	for rows.Next() {
		var memo Memo
		err := rows.Scan(&memo.ID, &memo.Title, &memo.Body, &memo.Tags)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

// MARK: - Delete()
func (i *itemRepository) Delete(ctx context.Context, memo *Memo) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	title := memo.Title
	body := memo.Body

	query := `DELETE FROM memos WHERE title = ? AND body = ?`
	result, err := tx.ExecContext(ctx, query, title, body)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("memo not exist") // カスタムエラーを返す
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// MARK: - SearchByKeyword()
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
func (i *itemRepository) SearchByTags(ctx context.Context, tags []string) ([]Memo, error) {
	if len(tags) == 0 {
		return []Memo{}, nil
	}

	// tagsスライスと同じ長さの文字列スライスを作成
	subqueries := make([]string, len(tags))
	// tagsスライスと同じ長さのインターフェーススライスを作成
	// SQLのプレースホルダの?に渡す引数を格納
	args := make([]interface{}, len(tags))
	for i, tag := range tags {
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
	if err = rows.Err(); err != nil {
		slog.Error("rows iteration failed", "error", err)
		return nil, err
	}
	// デバッグ用
	// slog.Info("query completed", "result_count", len(memos))
	return memos, nil
}
