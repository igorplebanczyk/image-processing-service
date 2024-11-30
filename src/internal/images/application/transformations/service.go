package transformations

import (
	"fmt"
	"image"
	"image-processing-service/src/internal/images/domain"
	"sync"
)

const (
	workerCount = 10
	queueSize   = 100
)

type Service struct {
	workerCoordinator *workerCoordinator
}

func NewService() *Service {
	return &Service{
		workerCoordinator: newWorkerCoordinator(workerCount, queueSize),
	}
}

func (s *Service) CreatePreview(bytes []byte) ([]byte, error) {
	return s.Apply(bytes, []domain.Transformation{
		{
			Type: domain.Resize,
			Options: map[domain.TransformationOptionType]float64{
				domain.Width:  200,
				domain.Height: 200,
			},
		},
	})
}

func (s *Service) Wait() {
	s.workerCoordinator.wait()
}

func (s *Service) Apply(imageBytes []byte, transformations []domain.Transformation) ([]byte, error) {
	packet, err := assemble(imageBytes, transformations)
	if err != nil {
		return nil, err
	}

	s.workerCoordinator.process(packet)

	return deassemble(packet)
}

type transformationPacket struct {
	img             image.Image
	format          string
	transformations []domain.Transformation
	responseChan    chan image.Image
	errChan         chan error
}

func assemble(imageBytes []byte, transformations []domain.Transformation) (*transformationPacket, error) {
	img, format, err := deserialize(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize image: %w", err)
	}

	return &transformationPacket{
		img:             img,
		format:          format,
		transformations: transformations,
		responseChan:    make(chan image.Image, 1),
		errChan:         make(chan error, 1),
	}, nil
}

func deassemble(packet *transformationPacket) ([]byte, error) {
	select {
	case resultImg := <-packet.responseChan:
		return serialize(resultImg, packet.format)
	case err := <-packet.errChan:
		return nil, err
	}
}

type workerCoordinator struct {
	workers  []*worker
	jobQueue chan *transformationPacket
	wg       sync.WaitGroup
}

func newWorkerCoordinator(workerCount, queueSize int) *workerCoordinator {
	o := &workerCoordinator{
		workers:  make([]*worker, workerCount),
		jobQueue: make(chan *transformationPacket, queueSize),
	}

	for i := 0; i < workerCount; i++ {
		w := newWorker(o.jobQueue)
		o.workers[i] = w
		go w.start()
	}
	return o
}

func (c *workerCoordinator) process(packet *transformationPacket) {
	c.wg.Add(1)
	c.jobQueue <- packet
}

func (c *workerCoordinator) wait() {
	c.wg.Wait()
}

type worker struct {
	jobQueue chan *transformationPacket
}

func newWorker(jobQueue chan *transformationPacket) *worker {
	return &worker{jobQueue: jobQueue}
}

func (w *worker) start() {
	for packet := range w.jobQueue {
		err := applyTransformations(packet)
		if err != nil {
			packet.errChan <- err
		} else {
			packet.responseChan <- packet.img
		}
		close(packet.responseChan)
		close(packet.errChan)
	}
}

func applyTransformations(packet *transformationPacket) error {
	var err error
	for _, t := range packet.transformations {
		switch t.Type {
		case domain.Resize:
			packet.img, err = resize(packet.img, t.Options)
		case domain.Crop:
			packet.img, err = crop(packet.img, t.Options)
		case domain.Rotate:
			packet.img, err = rotate(packet.img, t.Options)
		case domain.Grayscale:
			packet.img = grayscale(packet.img)
		case domain.Sepia:
			packet.img = sepia(packet.img)
		case domain.Invert:
			packet.img = invert(packet.img)
		case domain.AdjustBrightness:
			packet.img, err = adjustBrightness(packet.img, t.Options)
		case domain.AdjustContrast:
			packet.img, err = adjustContrast(packet.img, t.Options)
		case domain.AdjustSaturation:
			packet.img, err = adjustSaturation(packet.img, t.Options)
		case domain.Blur:
			packet.img, err = blur(packet.img, t.Options)
		case domain.Sharpen:
			packet.img, err = sharpen(packet.img, t.Options)
		default:
			err = fmt.Errorf("unsupported transformation type: %v", t.Type)
		}

		if err != nil {
			return fmt.Errorf("error applying transformation %v: %w", t.Type, err)
		}
	}
	return nil
}
