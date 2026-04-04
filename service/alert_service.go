package service

import (
	"context"

	"sentinelx/configs" // ✅ ADD THIS
	"sentinelx/models"
	"sentinelx/repository"
)

type AlertService struct {
	Repo *repository.AlertRepository
}

func NewAlertService(repo *repository.AlertRepository) *AlertService {
	return &AlertService{Repo: repo}
}

// 🔥 FIXED VERSION
func (s *AlertService) ProcessAlert(ctx context.Context, alert models.Alert) error {

	// ✅ ADD THIS (STEP 7)
	configs.Log("INFO", "processing alert", map[string]interface{}{
		"event_type": alert.Type,
		"source_ip":  alert.SourceIP,
		"tenant_id":  alert.TenantID,
		"alert_id":   alert.ID,
	})

	// save alert
	err := s.Repo.Save(ctx, alert)
	if err != nil {
		return err
	}

	return nil
}