type Memo struct {
	ID    int    `db:"id" json:"title"`
	Title string `db:"title" json:"title"`
	Body  string `db:"body" json:"body"`
	Image string `db:"image" json:"image"`
}

type ItemRepository interface {
	AddMemo(ctx context.Context, memo *Memo) error
}