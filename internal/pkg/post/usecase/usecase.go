package usecase

import (
	"github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/post"
)

type PostUseCase struct {
	PostRepository post.Repository
}

func NewPostUseCase(postRepository post.Repository) *PostUseCase {
	return &PostUseCase{
		PostRepository: postRepository,
	}
}
