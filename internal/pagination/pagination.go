package pagination

type Pagination struct {
	Limit  int
	Cursor int64
}

func (p *Pagination) Normalize() {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}
