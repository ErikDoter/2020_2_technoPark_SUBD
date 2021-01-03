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

type Threads []Thread
