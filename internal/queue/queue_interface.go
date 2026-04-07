package queue

type JobFunc func() error

type Priority int

const (
	High Priority = iota
	Medium
	Low
)

type JobQueue interface {
	Add(priority Priority, job JobFunc)
}
