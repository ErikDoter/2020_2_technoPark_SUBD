package repository

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/jackc/pgx"
	"time"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (r *ForumRepository) Create(title string, user string, slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	u := models.User{}
	err1 := r.db.QueryRow("select nickname from users where nickname = $1", user).
		Scan(&u.Nickname)
	if err1 != nil {
		return nil, &models.Error{Message: "can't find user with this nickname"}
	}
	_, err := r.db.Exec("insert into forums(slug, title, userf) values($1, $2, $3)", slug, title, u.Nickname)
	if err != nil {
		 err1 = r.db.QueryRow("select posts, slug, title, userf, threads from forums where slug = $1", slug).
			Scan(&forum.Posts, &forum.Slug, &forum.Title, &forum.User, &forum.Threads)
		 return &forum, &models.Error{Message: "exist"}
	}
	err = r.db.QueryRow("select posts, slug, threads, title, userf from forums where slug = $1", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	return &forum, nil
}

func (r *ForumRepository) Find(slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	err := r.db.QueryRow("select posts, slug, threads, title, userf from forums where slug = $1", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	return &forum, nil
}

func (r *ForumRepository) FindUsers(slug string, since string, desc bool, limit int) (*models.Users, *models.Error) {
	var query *pgx.Rows
	var err error
	forum := models.Forum{}
	err = r.db.QueryRow("select posts, slug, threads, title, userf from forums where slug = $1", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	user := models.User{}
	users := models.Users{}
	if desc {
		if since == "." {
			query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = $1        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = $2    ) as T ORDER BY  lower(T.nickname) DESC LIMIT $3", slug, slug, limit)
		} else {
			query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = $1        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = $2    ) as T where   lower(T.nickname) <   lower($3::text) ORDER BY lower(T.nickname) DESC LIMIT $4", slug, slug, since, limit)
		}
	} else {
		if since == "." {
			query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = $1        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = $2    ) as T ORDER BY  lower(T.nickname) LIMIT $3", slug, slug,  limit)
			if err != nil {
			}
		} else {
			query, err = r.db.Query("select T.about, T.email, T.fullname, T.nickname from (        SELECT u.about, u.email, u.fullname, u.nickname, u.id        from forums f join threads t on f.slug = t.forum        join users u on t.author = u.nickname        where f.slug = $1        union        SELECT uu.about, uu.email, uu.fullname, uu.nickname, uu.id        from forums ff join threads tt on ff.slug = tt.forum        join posts pp on pp.thread = tt.id        join users uu on uu.nickname = pp.author        where ff.slug = $2   ) as T where   lower(T.nickname) >   lower($3::text) ORDER BY  lower(T.nickname) LIMIT $4", slug, slug, since, limit)
		}
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
	var err error
	var i int
	if slugThread != "" {
		err = r.db.QueryRow("select id, title, author, forum, message, votes, created, slug from threads where slug = $1", slugThread).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug)
		if err == nil {
			return &thread, &models.Error{
				Message: "error",
			}
		}
	}
	err = r.db.QueryRow("select slug from forums where slug = $1", slug).
		Scan(&forum.Slug)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find",
		}
	}
	err = r.db.QueryRow("select nickname from users where nickname = $1", author).Scan(&user.Nickname)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find",
		}
	}
	date := time.Time{}
	if created != date {
		err = r.db.QueryRow("insert into threads(author, message, title, forum, slug, created) values($1, $2, $3, $4, $5, $6) returning id;", author, message, title, slug, slugThread, created).
			Scan(&i)
		r.db.QueryRow("select id, title, author, forum, message, votes, created, slug from threads where id = $1", i).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug)
		if err != nil {
			return &thread, &models.Error{
				Message: "error",
			}
		}
		return &thread, nil
	} else {
		err = r.db.QueryRow("insert into threads(author, message, title, forum, slug) values($1, $2, $3, $4, $5) returning id;", author, message, title, slug, slugThread).
			Scan(&i)
		r.db.QueryRow("select id, title, author, forum, message, votes, slug from threads where id = $1", i).
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
	err := r.db.QueryRow("select slug from forums where slug = $1", slug).
		Scan(&forum.Slug)
	if err != nil {
		return nil, &models.Error{
			Message: "can't find forum with this slug",
		}
	}
	var query *pgx.Rows
	if desc {
		if since == "" {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = $1 order by created DESC LIMIT $2", slug, limit)
		} else {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = $1  and created <= $2 order by created DESC LIMIT $3", slug, since, limit)
		}
	} else {
		if since == "" {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = $1 order by created LIMIT $2", slug, limit)
		} else {
			query, _ = r.db.Query("select id, title, author, forum, message, votes, slug, created from threads where forum = $1 and created >= $2 order by created LIMIT $3", slug, since, limit)
		}
	}
	for query.Next(){
		query.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		threads = append(threads, thread)
	}
	defer query.Close()
	return &threads, nil
}

