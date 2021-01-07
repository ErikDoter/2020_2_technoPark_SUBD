package repository

import (
	"database/sql"
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Check(id int) bool {
	post := models.Post{}
	var err error
	err = r.db.QueryRow("select id from posts where id = ?", id).
		Scan(&post.Id)
	if err != nil {
		return false
	}
	return true
}

func (r *PostRepository) DetailsUser(id int) models.User {
	user := models.User{}
	r.db.QueryRow("select u.about, u.email, u.fullname, u.nickname from posts p join users u on (p.id = ? and p.author = u.nickname)", id).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	return user
}

func (r *PostRepository) DetailsThread(id int) models.Thread {
	thread := models.Thread{}
	err := r.db.QueryRow("select t.author, t.slug, t.created, t.forum, t.id, t.message, t.title, t.votes from posts p join threads t on (p.id = ? and p.thread = t.id)", id).
		Scan(&thread.Author, &thread.Slug, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Title, &thread.Votes)
	if err != nil {
		fmt.Println(err)
	}
	return thread
}

func (r *PostRepository) DetailsForum(id int) models.Forum {
	forum := models.Forum{}
	r.db.QueryRow("select f.posts, f.slug, f.threads, f.title, f.user from posts p join forums f on (p.id = ? and p.forum = f.slug)", id).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	return forum
}

func (r *PostRepository) DetailsPost(id int) models.Post {
	post := models.Post{}
	r.db.QueryRow("select author, created, forum, id, message, isEdited, parent, thread from posts where id = ?", id).
		Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.Message, &post.IsEdited, &post.Parent, &post.Thread)
	return post
}

func (r *PostRepository) Update(id int, message string) models.Post {
	isEdited := true
	_, err := r.db.Exec("update posts set message = ?, isEdited = ? where id = ?", message, isEdited, id)
	if err != nil {
		fmt.Println(err)
	}
	post := models.Post{}
	r.db.QueryRow("select author, created, forum, id, message, isEdited, parent, thread from posts where id = ?", id).
		Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.Message, &post.IsEdited, &post.Parent, &post.Thread)
	return post
}

