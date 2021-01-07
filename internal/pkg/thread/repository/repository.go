package repository

import (
	"database/sql"
	"fmt"
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

func (r *ThreadRepository) CreatePosts(soi models.IdOrSlug, posts models.PostsMini) (*models.Posts, *models.Error) {
	thread := models.Thread{}
	post := models.Post{}
	postsAnswer := models.Posts{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = ?", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = ?", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	for _, value := range posts {
		if value.Parent != 0 {
			err = r.db.QueryRow("select thread from posts where id = ?", value.Parent).
				Scan(&post.Thread)
			if err != nil {
				return nil, &models.Error{Message: "can't find parent"}
			}
			if post.Thread != thread.Id {
				return nil, &models.Error{Message: "can't find parent"}
			}
		}
		_, err = r.db.Exec("insert into posts(author, message, parent, thread, forum) select ?, ?, ?, ?, forum from threads where id = ?", value.Author, value.Message, value.Parent, thread.Id, thread.Id)
		if err != nil {
			fmt.Println("bad")
		}
		err = r.db.QueryRow("select id, parent, author, message, isEdited, forum, thread, created from posts where id = (select MAX(id) from posts)").
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			fmt.Println(err)
		}
		postsAnswer = append(postsAnswer, post)
	}
	return &postsAnswer, nil
}


func (r *ThreadRepository) Update(soi models.IdOrSlug, title string, message string) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = ?", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = ?", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	_, err = r.db.Exec("update threads set message = ?, title = ? where id = ?;", message, title, thread.Id)
	id := thread.Id
	err = r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = ?", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		fmt.Println(err)
	}
	return &thread, nil
}

func (r *ThreadRepository) Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = ?", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = ?", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	_, err = r.db.Exec("insert into votes(nickname, thread, vote) value(?, ?, ?);", nickname, thread.Id, voice)
	if err != nil {
		fmt.Println(err)
		_, err = r.db.Exec("update votes set vote = ? where thread = ? and nickname = ?", voice, thread.Id, nickname)
	}
	id := thread.Id
	err = r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = ?", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		fmt.Println(err)
	}
	return &thread, nil
}

