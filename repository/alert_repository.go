package repository

import (
	"context"
	"sentinelx/models"
	"sentinelx/storage"
)

type AlertRepository struct{}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{}
}

// 🔥 wrapper over existing storage
func (r *AlertRepository) Save(ctx context.Context, alert models.Alert) {
	storage.SaveAlert(ctx, alert)
}