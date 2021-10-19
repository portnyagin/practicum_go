package service

import (
	"context"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/portnyagin/practicum_go/internal/app/model"
	"math"
	"sync"
)

type pool struct {
	sync.Mutex
	maxSize     int
	currentSize int
}

func newPool(maxSize int) *pool {
	var p pool
	p.maxSize = maxSize
	return &p
}

func (p *pool) Inc() {
	p.Lock()
	defer p.Unlock()
	p.currentSize += 1
}

func (p *pool) Dec() {
	p.Lock()
	defer p.Unlock()
	p.currentSize -= 1
}

func (p *pool) LessMax() bool {
	p.Lock()
	defer p.Unlock()
	return p.maxSize > p.currentSize
}

type deleteJob struct {
	UserID string
	chunk  []dto.BatchDeleteDTO
}

type DeleteService struct {
	// params
	pool         *pool
	taskSize     int
	dbRepository model.DeleteRepository
	jobChanel    chan deleteJob
}

func NewDeleteService(repoDB model.DeleteRepository, poolSize int, taskSize int) *DeleteService {
	var s DeleteService
	s.taskSize = taskSize
	s.pool = newPool(poolSize)
	s.dbRepository = repoDB
	s.jobChanel = make(chan deleteJob, poolSize)
	s.startWorkerPool()
	return &s
}

func split(batchSize int, src []dto.BatchDeleteDTO, resCh chan []dto.BatchDeleteDTO) {
	if batchSize <= 0 || len(src) == 0 {
		close(resCh)
		return
	}
	start := 0
	end := int(math.Min(float64(batchSize), float64(len(src)-start)))
	for start <= len(src) {
		resCh <- src[start:end]
		start = end + 1
		end = start + int(math.Min(float64(batchSize), float64(len(src)-start)))
	}
	close(resCh)
}

func (s *DeleteService) DeleteBatch(ctx context.Context, userID string, URLList []dto.BatchDeleteDTO) error {
	chanel := make(chan []dto.BatchDeleteDTO)
	go split(s.taskSize, URLList, chanel)
	for chunk := range chanel {
		s.jobChanel <- deleteJob{userID, chunk}
	}
	s.startWorkerPool()
	return nil
}

func (s *DeleteService) startWorkerPool() {
	for s.pool.LessMax() {
		s.pool.Inc()
		go func() {
			defer s.pool.Dec()
			for {
				for job := range s.jobChanel {
					err := s.dbRepository.BatchDelete(context.Background(), job.UserID, job.chunk)
					if err != nil {
						panic(err)
					}
					//time.Sleep(time.Second * 5)
				}
			}
		}()
	}
}
