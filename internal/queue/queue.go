package queue

type Job func() error
err := job()
if err != nil {
	for retry := 0; retry < 3; retry++ {
		err := job()

		if err == nil {
			break
		}

		log.Println("retry:", retry, "erro:", err)

		time.Sleep(1 * time.Second)
	}
}

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

				for retry := 0; retry < 3; retry++ {

					func() {
						defer func() {
							if r := recover(); r != nil {
								
							}
						}()

						job()
					}()

					break
				}
			}
		}(i)
	}
}

func (q *Queue) Add(job Job) {
	q.Jobs <- job
}