package transformation

import (
	"sync"
)

const (
	resizeTransform      = "resize"
	cropTransform        = "crop"
	rotateTransform      = "rotate"
	applyFilterTransform = "apply_filter"
	adjustTransform      = "adjust"
)

type Service struct {
	taskQueue   chan task
	workerCount int
	wg          sync.WaitGroup
}

func New(workerCount int, queueSize int) *Service {
	service := &Service{
		taskQueue:   make(chan task, queueSize),
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
	}
	service.startWorkerPool()
	return service
}

func (s *Service) Transform(imageBytes []byte, transformation string, options map[string]any) ([]byte, error) {
	task := task{
		imageBytes:     imageBytes,
		transformation: transformation,
		options:        options,
		result:         make(chan []byte, 1),
		err:            make(chan error, 1),
	}
	s.taskQueue <- task
	return <-task.result, <-task.err
}

func (s *Service) Close() {
	close(s.taskQueue)
	s.wg.Wait()
}
