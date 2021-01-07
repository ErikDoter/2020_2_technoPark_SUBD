package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/service"
)

type ServiceUseCase struct {
	ServiceRepository service.Repository
}

func NewServiceUseCase(serviceRepository service.Repository) *ServiceUseCase {
	return &ServiceUseCase{
		ServiceRepository: serviceRepository,
	}
}

func (u *ServiceUseCase) Status() *models.Status {
	status := u.ServiceRepository.Status()
	return status
}

func (u *ServiceUseCase) Clear() {
	u.ServiceRepository.Clear()
}
