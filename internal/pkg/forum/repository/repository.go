package repository

import (
	"database/sql"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"time"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (r *ForumRepository) Create(title string, user string, slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	u := models.User{}
	err1 := r.db.QueryRow("select nickname from users where nickname = ?", user).
		Scan(&u.Nickname)
	if err1 != nil {
		return nil, &models.Error{Message: "can't find user with this nickname"}
	}
	_, err := r.db.Exec("insert into forums(slug, title, user) value(?, ?, ?)", slug, title, u.Nickname)
	if err != nil {
		 err1 = r.db.QueryRow("select posts, slug, title, user, threads from forums where slug = ?", slug).
			Scan(&forum.Posts, &forum.Slug, &forum.Title, &forum.User, &forum.Threads)
		 return &forum, &models.Error{Message: "exist"}
	}
	err = r.db.QueryRow("select posts, slug, threads, title, user from forums where slug = ?", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	return &forum, nil
}

func (r *ForumRepository) Find(slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	err := r.db.QueryRow("select posts, slug, threads, title, user from forums where slug = ?", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	return &forum, nil
}

func (r *ForumRepository) FindUsers(slug string, since int, desc bool, limit int) (*models.Users, *models.Error) {
	var query *sql.Rows
	var err error
	forum := models.Forum{}
	err = r.db.QueryRow("select posts, slug, threads, title, user from forums where slug = ?", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	user := models.User{}
	users := models.Users{}
	if desc {
		query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = ?        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = ?    ) as T where T.id > ? ORDER BY lower(T.nickname) DESC LIMIT ?", slug, slug, since, limit)
	} else {
		query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = ?        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = ?    ) as T where T.id > ? ORDER BY lower(T.nickname) LIMIT ?", slug, slug, since, limit)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find forum with this slug"}
	}
	defer query.Close()
	for query.Next(){
		query.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		users = append(users, user)
	}
	return &users, nil
}

func (r *ForumRepository) CreateThread(slug string, title string, author string, message string, created time.Time, slugThread string) (*models.Thread, *models.Error) {

	forum := models.Forum{}
	user := models.User{}
	thread := models.Thread{}
	err := r.db.QueryRow("select slug from forums where slug = ?", slug).
		Scan(&forum.Slug)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find",
		}
	}
	err = r.db.QueryRow("select nickname from users where nickname = ?", author).Scan(&user.Nickname)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find",
		}
	}
	date := time.Time{}
	if created != date {
		_, err = r.db.Exec("insert into threads(author, message, title, forum, slug, created) value(?, ?, ?, ?, ?, ?);", author, message, title, slug, slugThread, created)
		r.db.QueryRow("select id, title, author, forum, message, votes, created, slug from threads where slug = ?", slugThread).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug)
		if err != nil {
			return &thread, &models.Error{
				Message: "error",
			}
		}
		return &thread, nil
	} else {
		_, err = r.db.Exec("insert into threads(author, message, title, forum, slug) value(?, ?, ?, ?, ?);", author, message, title, slug, slugThread)
		r.db.QueryRow("select id, title, author, forum, message, votes, slug from threads where slug = ?", slugThread).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug)
		if err != nil {
			return &thread, &models.Error{
				Message: "error",
			}
		}
		return &thread, nil
	}
}

func (r *ForumRepository) ShowThreads(slug string, limit int, since string, desc bool) (*models.Threads, *models.Error) {
	forum := models.Forum{}
	thread := models.Thread{}
	threads := models.Threads{}
	err := r.db.QueryRow("select slug from forums where slug = ?", slug).
		Scan(&forum.Slug)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	var query *sql.Rows
	if desc {
		if since == "" {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = ? order by created DESC LIMIT ?", slug, limit)
		} else {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = ?  and created <= ? order by created DESC LIMIT ?", slug, since, limit)
		}
	} else {
		if since == "" {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = ? order by created LIMIT ?", slug, limit)
		} else {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = ? and created >= ? order by created LIMIT ?", slug, since, limit)
		}
	}
	for query.Next(){
		query.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		threads = append(threads, thread)
	}
	defer query.Close()
	return &threads, nil
}

