package queue

import (
	"log"
	"time"
)

type PriorityQueue struct {
	high   chan JobFunc
	medium chan JobFunc
	low    chan JobFunc
}

func NewPriorityQueue(buffer int) *PriorityQueue {
	return &PriorityQueue{
		high:   make(chan JobFunc, buffer),
		medium: make(chan JobFunc, buffer),
		low:    make(chan JobFunc, buffer),
	}
}

func (q *PriorityQueue) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			log.Printf("Worker %d iniciado\n", workerID)

			for {
				select {
				case job := <-q.high:
					execute(job)

				case job := <-q.high:
					execute(job)

				case job := <-q.medium:
					execute(job)

				case job := <-q.low:
					execute(job)
				}
			}
		}(i)
	}
}

func (q *PriorityQueue) Add(priority Priority, job JobFunc) {
	switch priority {
	case High:
		q.high <- job
	case Medium:
		q.medium <- job
	default:
		q.low <- job
	}
}

func execute(job JobFunc) {
	for retry := 0; retry < 3; retry++ {
		err := func() (e error) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("panic recuperado:", r)
				}
			}()
			return job()
		}()

		if err == nil {
			return
		}

		log.Printf("Tentativa %d falhou: %v\n", retry+1, err)
		time.Sleep(1 * time.Second)
	}
}
