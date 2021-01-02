package repository

import (
	"database/sql"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByNickname(nickname string) (*models.User, *models.Error) {
	user := models.User{}
	err := r.db.QueryRow("select about, email, fullname, nickname from users where nickname = ?", nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, &models.Error{
			Message: "Can't find user with this nickname",
		}
	}

	return &user, nil
}
