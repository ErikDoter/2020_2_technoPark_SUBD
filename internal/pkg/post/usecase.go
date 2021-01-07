package post

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Usecase interface {
	Find(idStr string, related string) (*models.PostFull, *models.Error)
	Update(idStr string, message string) (*models.Post, *models.Error)
}
