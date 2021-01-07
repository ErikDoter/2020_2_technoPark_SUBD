package usecase

import (
	"fmt"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post"
	"strconv"
	"strings"
)

type PostUseCase struct {
	PostRepository post.Repository
}

func NewPostUseCase(postRepository post.Repository) *PostUseCase {
	return &PostUseCase{
		PostRepository: postRepository,
	}
}

func (u *PostUseCase) Find(idStr string, related string) (*models.PostFull, *models.Error) {
	id, _ := strconv.Atoi(idStr)
	var user = models.User{}
	var thread = models.Thread{}
	var forum = models.Forum{}
	var post = models.Post{}
	var postFull = models.PostFull{}
	if !u.PostRepository.Check(id) {
		return nil, &models.Error{Message: "can't find"}
	}
	if strings.Contains(related, "user") {
		user = u.PostRepository.DetailsUser(id)
	}
	if strings.Contains(related, "thread") {
		thread = u.PostRepository.DetailsThread(id)
	}
	if strings.Contains(related, "forum") {
		forum = u.PostRepository.DetailsForum(id)
	}
	post = u.PostRepository.DetailsPost(id)
	postFull.Post = post
	postFull.Author = user
	postFull.Thread = thread
	postFull.Forum = forum
	return &postFull, nil
}

func (u *PostUseCase) Update(idStr string, message string) (*models.Post, *models.Error){
	id, _ := strconv.Atoi(idStr)
	if !u.PostRepository.Check(id) {
		return nil, &models.Error{Message: "can't find"}
	}
	fmt.Println(id, message)
	post := u.PostRepository.Update(id, message)
	return &post, nil
}