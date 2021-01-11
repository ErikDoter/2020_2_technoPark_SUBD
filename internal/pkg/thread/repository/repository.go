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
	var forum string
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
	r.db.QueryRow("select forum from threads where id = $1", thread.Id).
		Scan(&forum)
	for _, value := range posts {
		err = r.db.QueryRow("insert into posts(author, message, parent, thread, created, forum) values($1, $2, $3, $4, $5, $6) RETURNING id, parent, author, message, isEdited,  forum, thread, created", value.Author, value.Message, value.Parent, thread.Id, now, forum).
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			if err.Error() == "ERROR: Not in this thread ID  (SQLSTATE P0001)" {
				return nil, &models.Error{Message: "can't find parent"}
			} else {
				return nil, &models.Error{Message: "can't find thread"}
			}
		}
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
	id := thread.Id
	if title != "" || message != "" {
		if title != "" {
			r.db.QueryRow("update threads set  title = $1 where id = $2 RETURNING id, title, author, forum, message, votes, slug, created", title, id).
				Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		}
		if message != "" {
			r.db.QueryRow("update threads set  message = $1 where id = $2 RETURNING id, title, author, forum, message, votes, slug, created", message, id).
				Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		}
	} else {
		r.db.QueryRow("select id, title, author, forum, message, votes, slug, created from threads where id = $1", id).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}
	return &thread, nil
}

func (r *ThreadRepository) Vote(soi models.IdOrSlug, nickname string, voice int ) (*models.Thread, *models.Error) {
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
	_, err = r.db.Exec("insert into votes(nickname, thread, vote) values($1, $2, $3) ON CONFLICT (thread, nickname) DO UPDATE SET vote = $3;", nickname, thread.Id, voice)
	if err != nil {
		return nil, &models.Error{Message: "can't find thread"}
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

func (r *ThreadRepository) CheckId(soi models.IdOrSlug) (int32, *models.Error) {
	var id int32
	var err error
	if soi.IsSlug {
		err = r.db.QueryRow("select id from threads where slug = $1", soi.Slug).
			Scan(&id)
	} else {
		err = r.db.QueryRow("select id from threads where id = $1", soi.Id).
			Scan(&id)
	}
	if err != nil {
		return 0, &models.Error{Message: "can't find"}
	}
	return id, nil
}


func (r *ThreadRepository) GetThreadPosts(threadID int32, desc bool, since string, limit int, sort string) (ps *models.Posts, err1 *models.Error) {
	posts := models.Posts{}
	query := ""

	var err error
	rows := &pgx.Rows{}
	if since != "" {
		switch sort {
		case "tree":
			query = "SELECT posts.id, posts.author, posts.forum, posts.thread, " +
				"posts.message, posts.parent, posts.isEdited, posts.created " +
				"FROM posts %s posts.thread = $1 ORDER BY posts.path[1] %s, posts.path %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path < p.path AND",
					"DESC",
					"DESC")
			} else {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path > p.path AND",
					"ASC",
					"ASC")
			}
		case "parent_tree":
			query =  "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created " +
				"FROM posts as p WHERE p.thread = $1 AND " +
				"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 %s %s %s"
			if desc {
				query = fmt.Sprintf(query, " AND p.path < (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] DESC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] DESC, p.path ")
			} else {
				query = fmt.Sprintf(query, " AND p.path > (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] ASC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] ASC, p.path ")
			}
		default:
			query = "SELECT id, author, forum, thread, message, parent, isEdited, created " +
				"FROM posts WHERE thread = $1 AND id %s $2 ORDER BY id %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "<", "DESC")
			} else {
				query = fmt.Sprintf(query, ">", "ASC")
			}
		}
		rows, err = r.db.Query(query, threadID, since, limit)
	} else {
		switch sort {
		case "tree":
			if desc {
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum, posts.thread, " +
					"posts.message, posts.parent, posts.isEdited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] DESC, posts.path DESC LIMIT $2")
			} else {
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum, posts.thread, " +
					"posts.message, posts.parent, posts.isEdited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] ASC, posts.path ASC LIMIT $2")
			}
		case "parent_tree":
			if desc {
				query = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0" +
					"ORDER BY p.path[1] DESC, p.path LIMIT $2)) " +
					"ORDER BY p.path[1] DESC, p.path"
			} else {
				query ="SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 " +
					"ORDER BY p.path[1] ASC, p.path LIMIT $2)) ORDER BY p.path[1] ASC, p.path"
			}
		default:
			if desc {
				query = "SELECT id, author, forum, thread, message, parent, isEdited, created " +
					"FROM posts WHERE thread = $1  ORDER BY id DESC LIMIT $2"
			} else {
				query = "SELECT id, author, forum, thread, message, parent, isEdited, created " +
					"FROM posts WHERE thread = $1 ORDER BY id ASC LIMIT $2"
			}
		}
		rows, err = r.db.Query(query, threadID, limit)
	}

	if err != nil {
		return &posts, &models.Error{Message: err.Error()}
	}

	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.Id, &p.Author, &p.Forum, &p.Thread, &p.Message, &p.Parent, &p.IsEdited, &p.Created)
		if err != nil {
			return &posts, &models.Error{Message: err.Error()}
		}
		posts = append(posts, *p)
	}
	return &posts, nil
}


