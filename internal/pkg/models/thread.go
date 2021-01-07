package models

type Thread struct {
	Author string `json:"author"`
	Created string `json:"created"`
	Forum string `json:"forum"`
	Id int32 `json:"id"`
	Message string `json:"message"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	Votes int32 `json:"votes"`
}

type IdOrSlug struct {
	Slug string
	Id int
	IsSlug bool
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title string `json:"title"`
}

type ThreadVote struct {
	Nickname string `json:"nickname"`
	Voice int `json:"voice"`
}

type Threads []Thread
