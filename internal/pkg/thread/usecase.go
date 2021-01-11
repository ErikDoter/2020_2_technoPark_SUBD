package thread

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Usecase interface {
	Find(slugOrId string) (*models.Thread, *models.Error)
	CreatePosts(slugOrId string, posts models.Posts) (*models.Posts, *models.Error)
	Update(slugOrId string, message string, title string) (*models.Thread, *models.Error)
	Vote(slugOrId string, nickname string, vote int) (*models.Thread, *models.Error)
	Posts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, *models.Error)
}
