package repository

import (
	"database/sql"
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
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
	_, err := r.db.Exec("truncate table posts")
	_, err = r.db.Exec("truncate table users")
	_, err = r.db.Exec("truncate table forums")
	_, err = r.db.Exec("truncate table threads")
	_, err = r.db.Exec("truncate table votes")
	if err != nil {
		fmt.Println(err)
	}

}