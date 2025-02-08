package workerpool

type Task func()

type WorkerPool struct {
	tasks chan Task
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		tasks: make(chan Task, 100),
	}

	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}

func (p *WorkerPool) worker() {
	for task := range p.tasks {
		task()
	}
}
