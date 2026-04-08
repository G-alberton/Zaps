package pagination

import "time"

type Response[T any] struct {
	Data       []T        `json:"data"`
	NextCursor *time.Time `json:"next_cursor,omitempty"`
	HasMore    bool       `json:"has_more"`
}
