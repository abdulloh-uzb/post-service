package repo

import (
	pbp "post-service/genproto/post"
)

type PostStorageI interface {
	Create(*pbp.PostReq) (*pbp.Post, error)
	DeletePost(id int) (*pbp.Empty, error)
	GetPost(id int) (*pbp.GetPostResponse, error)
	ListPost() (*pbp.Posts, error)
	UpdatePost(*pbp.Post) (*pbp.Post, error)
	GetPostByCustomerId(id int) (*pbp.Posts, error)
}
