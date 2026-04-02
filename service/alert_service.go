package service

import (
	"context"

	"sentinelx/models"
	"sentinelx/repository"
)

type AlertService struct {
	Repo *repository.AlertRepository
}

func NewAlertService(repo *repository.AlertRepository) *AlertService {
	return &AlertService{Repo: repo}
}

// 🔥 TEMP SAFE VERSION
func (s *AlertService) ProcessAlert(ctx context.Context, alert models.Alert) {

	// For now: directly save
	// (we will plug detection later safely)

	s.Repo.Save(ctx, alert)
}