package thread

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type Repository interface {
	Find(soi models.IdOrSlug) (*models.Thread, *models.Error)
	CreatePosts(soi models.IdOrSlug, posts models.Posts) (*models.Posts, *models.Error)
	Update(soi models.IdOrSlug, title string, message string) (*models.Thread, *models.Error)
	Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error)
	Check(soi models.IdOrSlug) *models.Error
	CheckId(soi models.IdOrSlug) (int32, *models.Error)
	GetThreadPosts(threadID int32, desc bool, since string, limit int, sort string) (ps *models.Posts, err1 *models.Error)
}
