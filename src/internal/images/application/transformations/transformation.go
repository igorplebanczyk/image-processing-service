package transformations

import (
	"image-processing-service/src/internal/images/domain"
	"sync"
)

const (
	workerCount = 10
	queueSize   = 100
)

type Service struct {
	taskQueue   chan task
	workerCount int
	wg          sync.WaitGroup
}

func New() *Service {
	service := &Service{
		taskQueue:   make(chan task, queueSize),
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
	}
	service.startWorkerPool()
	return service
}

func (s *Service) Apply(imageBytes []byte, transformation domain.Transformation) ([]byte, error) {
	task := task{
		imageBytes:     imageBytes,
		transformation: transformation.Type,
		options:        transformation.Options,
		result:         make(chan []byte, 1),
		err:            make(chan error, 1),
	}
	s.taskQueue <- task
	return <-task.result, <-task.err
}

func (s *Service) Wait() {
	s.wg.Wait()
}
