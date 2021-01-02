package user

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type UseCase interface {
	FindByNickname(nickname string) (*models.User, *models.Error)
}