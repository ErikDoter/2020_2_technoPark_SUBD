package repository

import (
	"database/sql"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (r *ThreadRepository) Find(soi models.IdOrSlug) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select author, created, forum, id, message, slug, title, votes from threads where slug = ?", soi.Slug).
			Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = r.db.QueryRow("select author, created, forum, id, message, slug, title, votes from threads where id = ?", soi.Id).
			Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}
	if err != nil {
		return nil, &models.Error{
			Message: "can't find thread",
		}
	}
	return &thread, nil
}
