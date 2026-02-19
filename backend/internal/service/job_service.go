package service

import (
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/repository"

	"github.com/sirupsen/logrus"
)

type JobService interface {
	CreateJob(job *domain.Job) error
	GetJob(jobID string) (*domain.Job, error)
	UpdateJob(job *domain.Job) error
	ProcessJob(jobID string) error
	GetAllJobs() ([]*domain.Job, error)
}

type JobServiceImpl struct {
	jobRepo           repository.JobRepository
	conversionService ConversionService
	logger            *logrus.Logger
}

func NewJobService(
	jobRepo repository.JobRepository,
	conversionService ConversionService,
	logger *logrus.Logger,
) *JobServiceImpl {
	return &JobServiceImpl{
		jobRepo:           jobRepo,
		conversionService: conversionService,
		logger:            logger,
	}
}

func (s *JobServiceImpl) CreateJob(job *domain.Job) error {
	s.logger.WithField("jobID", job.ID).Info("erstelle neuen Job")

	if err := s.jobRepo.Create(job); err != nil {
		s.logger.WithError(err).Error("fehler beim Erstellen des Jobs")
		return err
	}

	s.logger.WithField("jobID", job.ID).Info("Job erfolgreich erstellt")
	return nil
}

func (s *JobServiceImpl) GetJob(jobID string) (*domain.Job, error) {
	return s.jobRepo.FindByID(jobID)
}

func (s *JobServiceImpl) UpdateJob(job *domain.Job) error {
	return s.jobRepo.Update(job)
}

func (s *JobServiceImpl) ProcessJob(jobID string) error {
	s.logger.WithField("jobID", jobID).Info("starte Job-Verarbeitung")

	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		s.logger.WithError(err).Error("Job nicht gefunden")
		return err
	}

	job.UpdateStatus(domain.JobStatusProcessing)
	if err := s.jobRepo.Update(job); err != nil {
		s.logger.WithError(err).Error("fehler beim Aktualisieren des Job-Status")
		return err
	}

	if err := s.conversionService.Convert(job); err != nil {
		s.logger.WithError(err).Error("konvertierung fehlgeschlagen")
		job.SetError(err)
		s.jobRepo.Update(job)
		return err
	}

	job.UpdateStatus(domain.JobStatusCompleted)
	job.UpdateProgress(100)
	if err := s.jobRepo.Update(job); err != nil {
		s.logger.WithError(err).Error("fehler beim Aktualisieren des Job-Status")
		return err
	}

	s.logger.WithField("jobID", jobID).Info("Job-Verarbeitung erfolgreich abgeschlossen")
	return nil
}

func (s *JobServiceImpl) GetAllJobs() ([]*domain.Job, error) {
	return s.jobRepo.FindAll()
}
