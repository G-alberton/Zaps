package queue

import (
	"log"
	"time"
)

type Job func() error

type Queue struct {
	Jobs chan Job
}

func NewQueue(buffer int) *Queue {
	return &Queue{
		Jobs: make(chan Job, buffer),
	}
}

func (q *Queue) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for job := range q.Jobs {

				var err error

				for retry := 0; retry < 3; retry++ {

					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Println("panic recuperado:", r)
								err = nil
							}
						}()

						err = job()
					}()

					if err == nil {
						break
					}

					log.Println("retry:", retry, "erro:", err)
					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}
}

func (q *Queue) Add(job Job) {
	q.Jobs <- job
}
