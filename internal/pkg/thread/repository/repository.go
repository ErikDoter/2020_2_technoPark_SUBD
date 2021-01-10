package repository

import (
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/jackc/pgx"
	"time"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (r *ThreadRepository) Find(soi models.IdOrSlug) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select author, created, forum, id, message, slug, title, votes from threads where slug = $1", soi.Slug).
			Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = r.db.QueryRow("select author, created, forum, id, message, slug, title, votes from threads where id = $1", soi.Id).
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
	var now = time.Now()
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = $1", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = $1", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	for _, value := range posts {
		var i int
		err = r.db.QueryRow("select id from users where nickname = $1", value.Author).Scan(&i)
		if err != nil {
			return nil, &models.Error{Message: "can't find thread"}
		}
		if value.Parent != 0 {
			err = r.db.QueryRow("select thread from posts where id = $1", value.Parent).
				Scan(&post.Thread)
			if err != nil {
				return nil, &models.Error{Message: "can't find parent"}
			}
			if post.Thread != thread.Id {
				return nil, &models.Error{Message: "can't find parent"}
			}
		}
		_, err = r.db.Exec("insert into posts(author, message, parent, thread, created, forum) select $1, $2, $3, $4, $5, forum from threads where id = $6", value.Author, value.Message, value.Parent, thread.Id, now, thread.Id)
		if err != nil {
		}
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
		err = r.db.QueryRow("select id from threads where slug = $1", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = $1", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	if title != "" {
		_, err = r.db.Exec("update threads set  title = $1 where id = $2", title, thread.Id)
	}
	if message != "" {
		_, err = r.db.Exec("update threads set  message = $1 where id = $2", message, thread.Id)
	}
	id := thread.Id
	err = r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = $1", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return &thread, nil
}

func (r *ThreadRepository) Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var err error
	var i int
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = $1", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = $1", soi.Id).
			Scan(&thread.Id)
	}
	err2 := r.db.QueryRow("select id from users where nickname = $1", nickname).Scan(&i)
	if err2 != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
	}
	_, err = r.db.Exec("insert into votes(nickname, thread, vote) values($1, $2, $3);", nickname, thread.Id, voice)
	if err != nil {
		_, err = r.db.Exec("update votes set vote = $1 where thread = $2 and nickname = $3", voice, thread.Id, nickname)
	}
	id := thread.Id
	err = r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = $1", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return &thread, nil
}

func (r *ThreadRepository) Check(soi models.IdOrSlug) *models.Error {
	thread := models.Thread{}
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = $1", soi.Slug).
			Scan(&thread.Id)
	} else {
		err = r.db.QueryRow("select id from threads where id = $1", soi.Id).
			Scan(&thread.Id)
	}
	if err != nil {
		return &models.Error{Message: "can't find"}
	}
	return nil
}

func (r *ThreadRepository) PostsFlat(soi models.IdOrSlug, limit int, since int, desc bool) models.Posts {
	var query *pgx.Rows
	post := models.Post{}
	posts := models.Posts{}
	fmt.Println(limit)
	if soi.IsSlug {
		if desc {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread) where p.id < $2 order by p.id desc limit $3", soi.Slug, since, limit)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread) where p.id > $2 order by p.id  limit $3", soi.Slug, since, limit)
		}
	} else {
		if desc {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread) where p.id < $2 order by p.id desc limit $3", soi.Id, since, limit)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread) where p.id > $2 order by p.id  limit $3", soi.Id, since, limit)
		}
	}
	for query.Next(){
		err := query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, post)
	}
	defer query.Close()
	return posts
}

func (r *ThreadRepository) PostsParentTree (soi models.IdOrSlug, limit int, since int, desc bool) models.Posts {
	posts := models.Posts{}
	var query *pgx.Rows
	if soi.IsSlug {
		if desc {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) where p.id < $2 order by p.id desc limit $3", soi.Slug, since, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) where p.id < $2 order by p.id desc", soi.Slug, since)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		} else {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) order by p.id limit $2", soi.Slug, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) order by p.id", soi.Slug)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		}
	} else {
		if desc {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) where p.id < $2 order by p.id desc limit $3", soi.Id, since, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) where p.id < $2 order by p.id desc", soi.Id, since)
			}
			r.RecursiveParentTree(query, &posts)
			defer query.Close()
		} else {
			if since == 0 || since == 100000000 {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) order by p.id limit $2", soi.Id, limit)
			} else {
				query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) order by p.id", soi.Id)
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
	var query *pgx.Rows
	if soi.IsSlug {
		if since == 100000000 {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) where p.id > $2 order by p.id", soi.Slug, -since)
			r.RecursiveTree(query, &posts, limit, 0, desc)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.slug = $1 and t.id = p.thread and p.parent = 0) order by p.id", soi.Slug)
			if since > 0  {
				r.RecursiveTreeWithoutLimit(query, &posts, limit, since)
			} else {
				r.RecursiveTree(query, &posts, limit, since, desc)
			}
		}
		defer query.Close()
	} else {
		if since == 100000000 {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) where p.id > $2 order by p.id", soi.Id, -since)
			r.RecursiveTree(query, &posts, limit, 0, desc)
		} else {
			query, _ = r.db.Query("select p.author, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.created from threads t join posts p on (t.id = $1 and t.id = p.thread and p.parent = 0) order by p.id", soi.Id)
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

func (r * ThreadRepository) RecursiveTreeWithoutLimit(query *pgx.Rows, posts *models.Posts, limit int, since int) {
	post := models.Post{}
	var query1 *pgx.Rows
	for query.Next() {
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if post.Author == "" {
			return
		}
		*posts = append(*posts, post)
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = $1 order by id", post.Id)
		r.RecursiveTreeWithoutLimit(query1, posts, limit, since)
		query1.Close()
	}
}

func (r *ThreadRepository) RecursiveTree(query *pgx.Rows, posts *models.Posts, limit int, since int, desc bool) {
	post := models.Post{}
	var query1 *pgx.Rows
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
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = $1 order by id", post.Id)
		r.RecursiveTree(query1, posts, limit, since, desc)
		query1.Close()
	}
}

func (r *ThreadRepository) RecursiveParentTree(query *pgx.Rows, posts *models.Posts) {
	post := models.Post{}
	var query1 *pgx.Rows
	for query.Next() {
		query.Scan(&post.Author, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Created)
		if post.Author == "" {
			return
		}
		*posts = append(*posts, post)
		query1, _ = r.db.Query("select author, forum, id, isEdited, message, parent, thread, created from posts where parent = $1 order by id", post.Id)
		r.RecursiveParentTree(query1, posts)
		query1.Close()
	}
}

