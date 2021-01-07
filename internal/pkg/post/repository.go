package post

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Repository interface {
	Check(id int) bool
	DetailsUser(id int) models.User
	DetailsThread(id int) models.Thread
	DetailsForum(id int) models.Forum
	DetailsPost(id int) models.Post
	Update(id int, message string) models.Post
}
