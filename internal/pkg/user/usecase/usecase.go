package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/user"
)

type UserUseCase struct {
	UserRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) FindByNickname(nickname string) (*models.User, *models.Error) {
	user, err := u.UserRepository.FindByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return user, nil
}