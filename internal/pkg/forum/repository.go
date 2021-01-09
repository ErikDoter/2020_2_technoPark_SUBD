package forum

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"time"
)

type Repository interface {
	Create(title string, user string, slug string) (*models.Forum, *models.Error)
	Find(slug string) (*models.Forum, *models.Error)
	FindUsers(slug string, since int, desc bool, limit int) (*models.Users, *models.Error)
	CreateThread(slug string, title string, author string, message string, created time.Time, slugThread string) (*models.Thread, *models.Error)
	ShowThreads(slug string, limit int, since string, desc bool) (*models.Threads, *models.Error)
}
