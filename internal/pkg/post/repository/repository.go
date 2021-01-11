package repository

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/jackc/pgx"
)

type PostRepository struct {
	db  *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Check(id int) bool {
	post := models.Post{}
	var err error
	err = r.db.QueryRow("select id from posts where id = $1", id).
		Scan(&post.Id)
	if err != nil {
		return false
	}
	return true
}

func (r *PostRepository) DetailsUser(id int) models.User {
	user := models.User{}
	r.db.QueryRow("select u.about, u.email, u.fullname, u.nickname from posts p join users u on (p.id = $1 and p.author = u.nickname)", id).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	return user
}

func (r *PostRepository) DetailsThread(id int) models.Thread {
	thread := models.Thread{}
	err := r.db.QueryRow("select t.author, t.slug, t.created, t.forum, t.id, t.message, t.title, t.votes from posts p join threads t on (p.id = $1 and p.thread = t.id)", id).
		Scan(&thread.Author, &thread.Slug, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Title, &thread.Votes)
	if err != nil {
	}
	return thread
}

func (r *PostRepository) DetailsForum(id int) models.Forum {
	forum := models.Forum{}
	r.db.QueryRow("select f.posts, f.slug, f.threads, f.title, f.userf from posts p join forums f on (p.id = $1 and p.forum = f.slug)", id).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	return forum
}

func (r *PostRepository) DetailsPost(id int) models.Post {
	post := models.Post{}
	err := r.db.QueryRow("select author, forum, id, message, isEdited, parent, thread, created from posts where id = $1", id).
		Scan(&post.Author, &post.Forum, &post.Id, &post.Message, &post.IsEdited, &post.Parent, &post.Thread, &post.Created)
	if err != nil {
	}
	return post
}

func (r *PostRepository) Update(id int, message string) models.Post {
	isEdited := true
	var mes string
	post := models.Post{}
	r.db.QueryRow("select message from posts where id = $1", id).
		Scan(&mes)
	if message != "" && mes != message {
		r.db.QueryRow("update posts set message = $1, isEdited = $2 where id = $3 RETURNING author, forum, id, message, isEdited, parent, thread, created", message, isEdited, id).
			Scan(&post.Author, &post.Forum, &post.Id, &post.Message, &post.IsEdited, &post.Parent, &post.Thread, &post.Created)
	} else {
		r.db.QueryRow("select author, forum, id, message, isEdited, parent, thread, created from posts where id = $1", id).
			Scan(&post.Author, &post.Forum, &post.Id, &post.Message, &post.IsEdited, &post.Parent, &post.Thread, &post.Created)
	}
	return post
}

