package thread

import (
	"database/sql"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type Repository interface {
	Find(soi models.IdOrSlug) (*models.Thread, *models.Error)
	CreatePosts(soi models.IdOrSlug, posts models.PostsMini) (*models.Posts, *models.Error)
	Update(soi models.IdOrSlug, title string, message string) (*models.Thread, *models.Error)
	Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error)
	Check(soi models.IdOrSlug) *models.Error
	PostsFlat(soi models.IdOrSlug, limit int, since int, desc bool) models.Posts
	PostsTree (soi models.IdOrSlug, limit int, since int, desc bool) models.Posts
	RecursiveTree(query *sql.Rows, posts *models.Posts, limit int, since int, desc bool)
	PostsParentTree (soi models.IdOrSlug, limit int, since int, desc bool) models.Posts
	RecursiveParentTree(query *sql.Rows, posts *models.Posts)
	RecursiveTreeWithoutLimit(query *sql.Rows, posts *models.Posts, limit int, since int)
}
