package user

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type UseCase interface {
	FindByNickname(nickname string) (*models.User, *models.Error)
	Create(nickname string, fullname string, about string, email string) (*models.Users, *models.Error)
	Update(nickname string, fullname string, about string, email string) (*models.User, *models.Error)
}