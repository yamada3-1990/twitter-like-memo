-- メモのid, title, bodyを格納
CREATE TABLE IF NOT EXISTS memos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    body TEXT
);

-- タグのid, nameを格納
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- メモとタグを関連付ける(memosテーブルのidとtagsテーブルのidを対応させている)
CREATE TABLE IF NOT EXISTS memo_tags (
    memo_id INTEGER,
    tag_id INTEGER,
    FOREIGN KEY (memo_id) REFERENCES memos(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    PRIMARY KEY (memo_id, tag_id)
);               