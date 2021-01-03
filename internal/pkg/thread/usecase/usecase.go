package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread"
	"strconv"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
}

func NewThreadUseCase(threadRepository thread.Repository) *ThreadUseCase {
	return &ThreadUseCase{
		ThreadRepository: threadRepository,
	}
}

func (u *ThreadUseCase) Find(slugOrId string) (*models.Thread, *models.Error) {
	soi := models.IdOrSlug{}
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	thread, err := u.ThreadRepository.Find(soi)
	if err != nil {
		return nil, err
	}
	return thread, nil
}
