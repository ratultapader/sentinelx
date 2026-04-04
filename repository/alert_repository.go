package repository

import (
	"context"

	"sentinelx/configs"   // ✅ ADD
	"sentinelx/models"
	"sentinelx/storage"
)

type AlertRepository struct{}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{}
}

// 🔥 FIXED
func (r *AlertRepository) Save(ctx context.Context, alert models.Alert) error {

	err := storage.SaveAlert(ctx, alert)
	if err != nil {

		// ✅ ERROR LOGGING (STEP 8)
		configs.Log("ERROR", "failed to save alert", map[string]interface{}{
			"error":     err.Error(),
			"alert_id":  alert.ID,
			"source_ip": alert.SourceIP,
		})

		return err
	}

	return nil
}