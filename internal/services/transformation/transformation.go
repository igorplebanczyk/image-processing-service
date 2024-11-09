package transformation

import (
	"fmt"
	"image"
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

type task struct {
	imageBytes     []byte
	transformation string
	options        map[string]any
	result         chan []byte
	err            chan error
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

func (s *Service) startWorkerPool() {
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Service) worker() {
	defer s.wg.Done()
	for task := range s.taskQueue {
		result, err := s.processTask(task.imageBytes, task.transformation, task.options)
		task.result <- result
		task.err <- err
		close(task.result)
		close(task.err)
	}
}

func (s *Service) processTask(imageBytes []byte, transformation string, options map[string]any) ([]byte, error) {
	img, format, err := deserialize(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("error deserializing image: %w", err)
	}

	var transformed image.Image
	switch transformation {
	case resizeTransform:
		transformed, err = resize(img, options)
	case cropTransform:
		transformed, err = crop(img, options)
	case rotateTransform:
		transformed, err = rotate(img, options)
	case applyFilterTransform:
		transformed, err = applyFilter(img, options)
	case adjustTransform:
		transformed, err = adjust(img, options)
	default:
		return nil, fmt.Errorf("unknown transformation: %s", transformation)
	}
	if err != nil {
		return nil, fmt.Errorf("error during %s transformation: %w", transformation, err)
	}

	return serialize(transformed, format)
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
