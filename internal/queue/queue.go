package queue

import (
	"fmt"
	"log"
	"time"
)

type Queue struct {
	Jobs chan JobFunc
}

func NewQueue(buffer int) *Queue {
	return &Queue{
		Jobs: make(chan JobFunc, buffer),
	}
}

func (q *Queue) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func(workerID int) {

			log.Printf("Worker %d iniciado\n", workerID)

			for job := range q.Jobs {

				log.Printf("Worker %d processando job\n", workerID)

				var err error

				for retry := 0; retry < 3; retry++ {

					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Printf("Worker %d panic recuperado: %v\n", workerID, r)
								err = fmt.Errorf("panic: %v", r)
							}
						}()

						err = job()
					}()

					if err == nil {
						break
					}

					log.Printf("Worker %d retry %d erro: %v\n", workerID, retry+1, err)
					time.Sleep(time.Duration(retry+1) * time.Second)
				}

				if err != nil {
					log.Printf("Worker %d: job falhou após 3 tentativas: %v\n", workerID, err)
				}
			}

		}(i)
	}
}

func (q *Queue) Add(_ Priority, job JobFunc) {
	q.Jobs <- job
}
