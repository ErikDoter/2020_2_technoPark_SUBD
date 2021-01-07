package thread

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Repository interface {
	Find(soi models.IdOrSlug) (*models.Thread, *models.Error)
	CreatePosts(soi models.IdOrSlug, posts models.PostsMini) (*models.Posts, *models.Error)
	Update(soi models.IdOrSlug, title string, message string) (*models.Thread, *models.Error)
	Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error)
}
