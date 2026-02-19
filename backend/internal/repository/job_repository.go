package repository

import (
	"pptx2mp4/backend/internal/domain"
	"sync"
)

type JobRepository interface {
	Create(job *domain.Job) error
	FindByID(id string) (*domain.Job, error)
	Update(job *domain.Job) error
	Delete(id string) error
	FindAll() ([]*domain.Job, error)
}

type InMemoryJobRepository struct {
	jobs map[string]*domain.Job
	mu   sync.RWMutex
}

func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]*domain.Job),
	}
}

func (r *InMemoryJobRepository) Create(job *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.ID]; exists {
		return domain.ErrJobAlreadyExists
	}

	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryJobRepository) FindByID(id string) (*domain.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, domain.ErrJobNotFound
	}

	return job, nil
}

func (r *InMemoryJobRepository) Update(job *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.ID]; !exists {
		return domain.ErrJobNotFound
	}

	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryJobRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[id]; !exists {
		return domain.ErrJobNotFound
	}

	delete(r.jobs, id)
	return nil
}

func (r *InMemoryJobRepository) FindAll() ([]*domain.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]*domain.Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}
