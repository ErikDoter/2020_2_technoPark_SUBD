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

func (u *UserUseCase) Create(nickname string, fullname string, about string, email string) (*models.Users, *models.Error) {
	users, err := u.UserRepository.Create(nickname, fullname, about, email)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (u *UserUseCase) Update(nickname string, fullname string, about string, email string) (*models.User, *models.Error) {
	user, err := u.UserRepository.Update(nickname, fullname, about, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}