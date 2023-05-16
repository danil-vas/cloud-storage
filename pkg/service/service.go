package service

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/danil-vas/cloud-storage/pkg/repository"
)

type Authorization interface {
	CreateUser(user cloud_storage.User) (int, error)
	CreateMainDirectory(id int, login string) error
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type File interface {
	PathUploadFile(userId int, objectId int) (string, error)
}

type Service struct {
	Authorization
	File
}

func NewServices(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		File:          NewFileService(repos.File),
	}
}
