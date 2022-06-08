package storage

import (
	"Project/template-service/storage/postgres"
	"Project/template-service/storage/repo"

	"github.com/jmoiron/sqlx"
)

//IStorage ...
type IStorage interface {
	Post() repo.PostStorageI
}

type storagePg struct {
	db       *sqlx.DB
	PostRepo repo.PostStorageI
}

//NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		PostRepo: postgres.NewUserRepo(db),
	}
}

func (s storagePg) Post() repo.PostStorageI {
	return s.PostRepo
}

