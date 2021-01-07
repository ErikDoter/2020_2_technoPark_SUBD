package thread

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Usecase interface {
	Find(slugOrId string) (*models.Thread, *models.Error)
	CreatePosts(slugOrId string, posts models.PostsMini) (*models.Posts, *models.Error)
	Update(slugOrId string, message string, title string) (*models.Thread, *models.Error)
	Vote(slugOrId string, nickname string, vote int) (*models.Thread, *models.Error)
}
