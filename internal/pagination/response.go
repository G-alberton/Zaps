package pagination

type Response[T any] struct {
	Data       []T   `json:"data"`
	NextCursor int64 `json:"next_cursor,omitempty"`
	HasMore    bool  `json:"has_more"`
}
