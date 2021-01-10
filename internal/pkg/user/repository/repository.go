package repository

import (
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByNickname(nickname string) (*models.User, *models.Error) {
	user := models.User{}
	err := r.db.QueryRow("select about, email, fullname, nickname from users where nickname = $1", nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, &models.Error{
			Message: "Can't find user with this nickname",
		}
	}

	return &user, nil
}

func (r *UserRepository) Create(nickname string, fullname string, about string, email string) (*models.Users, *models.Error) {
	users := models.Users{}
	user := models.User{}
	_, err := r.db.Exec("insert into users(nickname, fullname, email, about) values($1, $2, $3, $4)", nickname, fullname, email, about)
	if err != nil {
		query, _ := r.db.Query("select about, email, fullname, nickname from users where nickname = $1 union select about, email, fullname, nickname from users where email = $2", nickname, email)
		defer query.Close()
		for query.Next(){
			query.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
			users = append(users, user)
		}
		return &users, &models.Error{Message: "409"}
	}
	user = models.User{
		About:    about,
		Email:    email,
		Fullname: fullname,
		Nickname: nickname,
	}
	users = append(users, user)
	return &users, nil
}

func (r *UserRepository) Update(nickname string, fullname string, about string, email string) (*models.User, *models.Error) {
	user := models.User{
		About:    about,
		Email:    email,
		Fullname: fullname,
		Nickname: nickname,
	}
	query, err := r.db.Query("select about, email, fullname, nickname from users where nickname = $1", nickname)
	defer query.Close()
	if !query.Next() {
		return nil, &models.Error{
			Message: "don't exist",
		}
	}
	sql := "Update users set"
	if about != "" {
		sql += fmt.Sprintf(" about = '%s',", about)
	}
	if email != "" {
		sql += fmt.Sprintf(" email = '%s',", email)
	}
	if fullname != "" {
		sql += fmt.Sprintf(" fullname = '%s',", fullname)
	}
	sql = sql[:len(sql) - 1]
	sql += " where nickname = $1"
	if about != "" || email != "" || fullname != "" {
		_, err = r.db.Exec(sql, nickname)
		if err != nil {
			return nil, &models.Error{Message: "conflict"}
		}
	}
	err = r.db.QueryRow("select about, email, fullname, nickname from users where nickname = $1", nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
			return nil, &models.Error{Message: "conflict"}
		}

	return &user, nil
}

