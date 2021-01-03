package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/forum"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	uuid2 "github.com/satori/go.uuid"
	"strconv"
)

type ForumUseCase struct {
	ForumRepository forum.Repository
}

func NewForumUseCase(forumRepository forum.Repository) *ForumUseCase {
	return &ForumUseCase{
		ForumRepository: forumRepository,
	}
}

func (u *ForumUseCase) Create(title string, user string, slug string) (*models.Forum, *models.Error) {
	forum, err := u.ForumRepository.Create(title, user, slug)
	return forum, err
}

func (u *ForumUseCase) Find(slug string) (*models.Forum, *models.Error) {
	forum, err := u.ForumRepository.Find(slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (u *ForumUseCase) FindUsers(slug string, since int, desc bool, limit int) (*models.Users, *models.Error) {
	users, err := u.ForumRepository.FindUsers(slug, since, desc, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *ForumUseCase) CreateThread(slug string, title string, author string, message string, created string, slugThread string) (*models.Thread, *models.Error) {
	if slugThread == "" {
		slugThread = uuid2.NewV4().String()
	}
	thread, err := u.ForumRepository.CreateThread(slug, title, author, message, created, slugThread)
	return thread, err
}

func (u *ForumUseCase) ShowThreads(slug string, limit string, since string, desc string) (*models.Threads, *models.Error) {
	var l, s int
	var d bool
	if limit == "" {
		l = 100
	} else {
		l, _ = strconv.Atoi(limit)
	}
	if desc == "" || desc == "false" {
		d = false
	} else {
		d = true
	}
	if since == "" {
		s = 0
	} else {
		s, _ = strconv.Atoi(since)
	}
	threads, err := u.ForumRepository.ShowThreads(slug, l, s, d)
	if err != nil {
		return nil, err
	}
	return threads, nil
}