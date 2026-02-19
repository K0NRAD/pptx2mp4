package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID             string            `json:"jobId"`
	Status         JobStatus         `json:"status"`
	Progress       int               `json:"progress"`
	Error          string            `json:"error,omitempty"`
	Config         *ConversionConfig `json:"config"`
	OriginalFile   string            `json:"originalFile"`
	OutputFile     string            `json:"outputFile,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	CompletedAt    *time.Time        `json:"completedAt,omitempty"`
}

func NewJob(originalFile string, config *ConversionConfig) *Job {
	now := time.Now()
	return &Job{
		ID:           uuid.New().String(),
		Status:       JobStatusPending,
		Progress:     0,
		Config:       config,
		OriginalFile: originalFile,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (j *Job) UpdateStatus(status JobStatus) {
	j.Status = status
	j.UpdatedAt = time.Now()

	if status == JobStatusCompleted || status == JobStatusFailed {
		now := time.Now()
		j.CompletedAt = &now
	}
}

func (j *Job) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	j.Progress = progress
	j.UpdatedAt = time.Now()
}

func (j *Job) SetError(err error) {
	j.Status = JobStatusFailed
	j.Error = err.Error()
	now := time.Now()
	j.CompletedAt = &now
	j.UpdatedAt = now
}

func (j *Job) SetOutputFile(outputFile string) {
	j.OutputFile = outputFile
	j.UpdatedAt = time.Now()
}

func (j *Job) IsCompleted() bool {
	return j.Status == JobStatusCompleted
}

func (j *Job) IsFailed() bool {
	return j.Status == JobStatusFailed
}

func (j *Job) IsProcessing() bool {
	return j.Status == JobStatusProcessing || j.Status == JobStatusPending
}
