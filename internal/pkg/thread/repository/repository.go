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
		err = r.db.QueryRow("select id, parent, author, message, isEdited, forum, thread, created from posts where id = (select MAX(id) from posts)").
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
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
		_, err = r.db.Exec("update votes set vote = ? where thread = ? and nickname = ?", voice, thread.Id, nickname)
	}
	id := thread.Id
	err = r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = ?", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return &thread, nil
}

func (r *ThreadRepository) Check(soi models.IdOrSlug) *models.Error {
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
		return &models.Error{Message: "can't find"}
	}
	return nil
}

func (r *ThreadRepository) PostsFlat(soi models.IdOrSlug, limit int, since int, desc bool) models.Posts {
	var query *sql.Rows
	post := models.Post{}
	posts := models.Posts{}
	if soi.IsSlug {
		if desc {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread) where p.id < ? order by p.id desc limit ?", soi.Slug, since, limit)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread) where p.id > ? order by p.id  limit ?", soi.Slug, since, limit)

		}
	} else {
		if desc {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread) where p.id < ? order by p.id desc limit ?", soi.Id, since, limit)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread) where p.id > ? order by p.id  limit ?", soi.Id, since, limit)
		}
	}
	for query.Next(){
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		posts = append(posts, post)
	}
	defer query.Close()
	return posts
}

func (r *ThreadRepository) PostsParentTree (soi models.IdOrSlug, limit int, since int, desc bool) models.Posts {
	posts := models.Posts{}
	var query *sql.Rows
	if soi.IsSlug {
		if desc {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) where p.id < ? order by p.id desc limit ?", soi.Slug, since, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) where p.id < ? order by p.id desc", soi.Slug, since)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		} else {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) order by p.id limit ?", soi.Slug, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) order by p.id", soi.Slug)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		}
	} else {
		if desc {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) where p.id < ? order by p.id desc limit ?", soi.Id, since, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) where p.id < ? order by p.id desc", soi.Id, since)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		} else {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) order by p.id limit ?", soi.Id, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) order by p.id", soi.Id)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		}
	}
	if since != 0 && since != 100000000 {
		posts = deleteSinceParent(posts, since)
	}
	return posts
}

func (r *ThreadRepository) PostsTree (soi models.IdOrSlug, limit int, since int, desc bool) models.Posts {
	posts := models.Posts{}
	var query *sql.Rows
	if soi.IsSlug {
		if since == 100000000 {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) where p.id > ? order by p.id", soi.Slug, -since)
			r.RecursiveTree(query, &posts, limit, 0, desc)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = ? and t.id = p.thread and p.parent = 0) order by p.id", soi.Slug)
			if since > 0  {
				r.RecursiveTreeWithoutLimit(query, &posts, limit, since)
			} else {
				r.RecursiveTree(query, &posts, limit, since, desc)
			}
		}
		defer query.Close()
	} else {
		if since == 100000000 {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) where p.id > ? order by p.id", soi.Id, -since)
			r.RecursiveTree(query, &posts, limit, 0, desc)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = ? and t.id = p.thread and p.parent = 0) order by p.id", soi.Id)
			if since > 0  {
				r.RecursiveTreeWithoutLimit(query, &posts, limit, since)
			} else {
				r.RecursiveTree(query, &posts, limit, since, desc)
			}
		}
		defer query.Close()
	}
	if desc {
		posts = reverse(posts)
		if since == 0 || since == 100000000 {
			posts = Limit(posts, limit)
		}
	}
	if since != 0 && since != 100000000 {
		posts = deleteSince(posts, since, limit)
	}
	return posts
}

func deleteSinceParent(posts models.Posts, since int) models.Posts {
	var i int
	for index, value := range posts {
		if value.Id == since {
			i = index + 1
			break
		}
	}
	if i <= len(posts) {
		posts = posts[i:]
	} else {
		posts = posts[i-1:]
	}
	return posts
}

func deleteSince(posts models.Posts, since int, limit int) models.Posts {
	var i int
	for index, value := range posts {
		if value.Id == since {
			i = index + 1
			break
		}
	}
	if i <= len(posts) {
		posts = posts[i:]
	} else {
		posts = posts[i-1:]
	}
	if limit <= len(posts) {
		posts = posts[:limit]
	}
	return posts
}

func reverse(posts models.Posts) models.Posts {
	for i := 0; i < len(posts)/2; i++ {
		j := len(posts) - i - 1
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts
}

func Limit(posts models.Posts, limit int) models.Posts {
	if limit <= len(posts) {
		posts = posts[:limit]
	}
	return posts
}

func (r * ThreadRepository) RecursiveTreeWithoutLimit(query *sql.Rows, posts *models.Posts, limit int, since int) {
	post := models.Post{}
	var query1 *sql.Rows
	for query.Next() {
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if post.Author == "" {
			return
		}
		*posts = append(*posts, post)
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = ?", post.Id)
		r.RecursiveTreeWithoutLimit(query1, posts, limit, since)
		query1.Close()
	}
}

func (r *ThreadRepository) RecursiveTree(query *sql.Rows, posts *models.Posts, limit int, since int, desc bool) {
	post := models.Post{}
	var query1 *sql.Rows
	if len(*posts) == limit && !desc {
		return
	}
	for query.Next() {
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if post.Author == "" {
			return
		}
		if len(*posts) == limit && !desc {
			return
		}
		*posts = append(*posts, post)
		if len(*posts) == limit && !desc {
			return
		}
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = ?", post.Id)
		r.RecursiveTree(query1, posts, limit, since, desc)
		query1.Close()
	}
}

func (r *ThreadRepository) RecursiveParentTree(query *sql.Rows, posts *models.Posts) {
	post := models.Post{}
	var query1 *sql.Rows
	for query.Next() {
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if post.Author == "" {
			return
		}
		*posts = append(*posts, post)
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = ?", post.Id)
		r.RecursiveParentTree(query1, posts)
		query1.Close()
	}
}

