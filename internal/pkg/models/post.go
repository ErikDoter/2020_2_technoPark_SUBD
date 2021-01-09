package models

import "time"

type Post struct {
	Author string `json:"author"`
	Created time.Time `json:"created"`
	Forum string `json:"forum"`
	Id int `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int64 `json:"parent"`
	Thread int32 `json:"thread"`
}

type PostFull struct {
	Author User `json:"author"`
	Forum Forum `json:"forum"`
	Post Post `json:"post"`
	Thread Thread `json:"thread"`
}

type PostMini struct {
	Parent int64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
}

type PostMessage struct {
	Message string `json:"message"`
}

type PostsMini []PostMini
type Posts []Post