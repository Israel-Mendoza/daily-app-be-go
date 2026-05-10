package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/repository"
)

type DailySessionService interface {
	FindAll(ctx context.Context) ([]db.DailySession, error)
	FindById(ctx context.Context, id int32) (*db.DailySession, error)
	Save(ctx context.Context, ds db.DailySession) (db.DailySession, error)
	DeleteById(ctx context.Context, id int32) error
}

type dailySessionService struct {
	repo repository.DailySessionRepository
}

func NewDailySessionService(repo repository.DailySessionRepository) DailySessionService {
	return &dailySessionService{repo: repo}
}

func (s *dailySessionService) FindAll(ctx context.Context) ([]db.DailySession, error) {
	return s.repo.FindAll(ctx)
}

func (s *dailySessionService) FindById(ctx context.Context, id int32) (*db.DailySession, error) {
	ds, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

func (s *dailySessionService) Save(ctx context.Context, ds db.DailySession) (db.DailySession, error) {
	return s.repo.Create(ctx, ds.SessionDate, ds.RawNotesBlob, ds.GeneratedScript)
}

func (s *dailySessionService) DeleteById(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}
