package repository

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/jmoiron/sqlx"
	"time"
)

type Authorization interface {
	CreateUser(user cloud_storage.User) (int, error)
	CreateMainDirectory(id int, login string) error
	GetUser(login, password string) (cloud_storage.User, error)
}

type File interface {
	PathUploadFile(userId int, objectId int) (string, error)
	AddUploadFileToUser(userId int, objectId int, originalName string, serverName string, size int64, create_time time.Time) (int, error)
	OriginalFileName(objectId int) (string, error)
	DeleteFile(userId int, objectId int) error
	GetAvailableMemory(userId int) (int, error)
	CheckAccessToObject(userId int, objectId int) (bool, error)
	GetTypeObject(objectId int) (string, error)
}

type Directory interface {
	AddDirectory(userId int, objectId int, nameDirectory string) error
	GetDirectoriesAndFiles(userId int, objectId int) ([]cloud_storage.Node, error)
	DeleteDirectory(objectId int) error
	GetIdMainDirectory(userID int) (int, error)
}

type User interface {
	GetUser(userId int) (cloud_storage.UserInfo, error)
}

type Repository struct {
	Authorization
	File
	Directory
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		File:          NewFilePostgres(db),
		Directory:     NewDirectoryPostgres(db),
		User:          NewUserPostgres(db),
	}
}
