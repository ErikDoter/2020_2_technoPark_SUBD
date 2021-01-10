package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/thread"
	"strconv"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
}

func NewThreadUseCase(threadRepository thread.Repository) *ThreadUseCase {
	return &ThreadUseCase{
		ThreadRepository: threadRepository,
	}
}

func (u *ThreadUseCase) Find(slugOrId string) (*models.Thread, *models.Error) {
	soi := models.IdOrSlug{}
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	thread, err := u.ThreadRepository.Find(soi)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (u *ThreadUseCase) CreatePosts(slugOrId string, posts models.PostsMini) (*models.Posts, *models.Error) {
	soi := models.IdOrSlug{}
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	postsAnswer, err := u.ThreadRepository.CreatePosts(soi, posts)
	if err != nil {
		return nil, err
	}
	return postsAnswer, nil
}

func (u *ThreadUseCase) Update(slugOrId string, message string, title string) (*models.Thread, *models.Error) {
	soi := models.IdOrSlug{}
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	thread, err := u.ThreadRepository.Update(soi, title, message)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (u *ThreadUseCase) Vote(slugOrId string, nickname string, vote int) (*models.Thread, *models.Error) {
	soi := models.IdOrSlug{}
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	thread, err := u.ThreadRepository.Vote(soi, nickname, vote)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (u *ThreadUseCase) Posts(slugOrId string, limit string, since string, sort string, desc string) (*models.Posts, *models.Error) {
	soi := models.IdOrSlug{}
	var l int
	var d bool
	result, err1 := strconv.Atoi(slugOrId)
	if err1 != nil {
		soi.Slug = slugOrId
		soi.IsSlug = true
	} else {
		soi.Id = result
		soi.IsSlug = false
	}
	id, err := u.ThreadRepository.CheckId(soi)
	if err != nil {
		return nil, err
	}
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
	posts, _ := u.ThreadRepository.GetThreadPosts(id, d, since, l, sort)
	return posts, nil
}