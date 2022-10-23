package storage

import (
	"post-service/storage/postgres"
	"post-service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Post() repo.PostStorageI
}

type storagePg struct {
	db       *sqlx.DB
	postRepo repo.PostStorageI
}

func NewStorage(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}

func (s *storagePg) Post() repo.PostStorageI {
	return s.postRepo
}
