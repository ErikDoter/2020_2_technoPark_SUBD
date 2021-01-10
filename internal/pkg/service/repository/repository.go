package repository

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/jackc/pgx"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (r *ServiceRepository) Status() *models.Status {
	status := models.Status{}
	r.db.QueryRow("select count(*) from users").
		Scan(&status.User)
	r.db.QueryRow("select count(*) from forums").
		Scan(&status.Forum)
	r.db.QueryRow("select count(*) from threads").
		Scan(&status.Thread)
	r.db.QueryRow("select count(*) from posts").
		Scan(&status.Post)
	return &status
}

func (r *ServiceRepository) Clear(){
	r.db.Exec("truncate table posts")
	r.db.Exec("truncate table users")
	r.db.Exec("truncate table forums")
	r.db.Exec("truncate table threads")
	r.db.Exec("truncate table votes")
}