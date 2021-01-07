package service

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Repository interface {
	Status() *models.Status
	Clear()
}
