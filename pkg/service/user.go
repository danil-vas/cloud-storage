package service

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/danil-vas/cloud-storage/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(userId int) (cloud_storage.UserInfo, error) {
	return s.repo.GetUser(userId)
}
