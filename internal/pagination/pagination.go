package pagination

import "time"

type Pagination struct {
	Limit  int
	Cursor time.Time
}

func (p *Pagination) Normalize() {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}

	if p.Cursor.IsZero() {

	}
}
