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
	orchestrator *Orchestrator
}

func NewService() *Service {
	return &Service{
		orchestrator: NewOrchestrator(workerCount, queueSize),
	}
}

type transformationPacket struct {
	img             image.Image
	transformations []domain.Transformation
	responseChan    chan image.Image
	errChan         chan error
}

func (s *Service) Apply(imageBytes []byte, transformations []domain.Transformation) ([]byte, error) {
	packet, err := assemble(imageBytes, transformations)
	if err != nil {
		return nil, err
	}

	s.orchestrator.Process(packet)

	select {
	case resultImg := <-packet.responseChan:
		return deassemble(resultImg)
	case err := <-packet.errChan:
		return nil, err
	}
}

func assemble(imageBytes []byte, transformations []domain.Transformation) (*transformationPacket, error) {
	img, _, err := deserialize(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize image: %w", err)
	}

	return &transformationPacket{
		img:             img,
		transformations: transformations,
		responseChan:    make(chan image.Image, 1),
		errChan:         make(chan error, 1),
	}, nil
}

func deassemble(img image.Image) ([]byte, error) {
	return serialize(img, "png")
}

type Orchestrator struct {
	workers  []*worker
	jobQueue chan *transformationPacket
	wg       sync.WaitGroup
}

func NewOrchestrator(workerCount, queueSize int) *Orchestrator {
	o := &Orchestrator{
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

func (o *Orchestrator) Process(packet *transformationPacket) {
	o.wg.Add(1)
	o.jobQueue <- packet
}

func (o *Orchestrator) Wait() {
	o.wg.Wait()
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
