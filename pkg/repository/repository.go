package repository

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user cloud_storage.User) (int, error)
	CreateMainDirectory(id int, login string) error
	GetUser(login, password string) (cloud_storage.User, error)
}

type File interface {
	PathUploadFile(userId int, objectId int) (string, error)
}

type Repository struct {
	Authorization
	File
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		File:          NewFilePostgres(db),
	}
}
