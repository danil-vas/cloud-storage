package service

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/danil-vas/cloud-storage/pkg/repository"
)

type DirectoryService struct {
	repo repository.Directory
}

func NewDirectoryService(repo repository.Directory) *DirectoryService {
	return &DirectoryService{repo: repo}
}

func (s *DirectoryService) AddDirectory(userId int, objectId int, nameDirectory string) (int, error) {
	return s.repo.AddDirectory(userId, objectId, nameDirectory)
}

func (s *DirectoryService) GetDirectoriesAndFiles(userId int, objectId int) ([]cloud_storage.Node, error) {
	return s.repo.GetDirectoriesAndFiles(userId, objectId)
}

func (s *DirectoryService) DeleteDirectory(objectId int) error {
	return s.repo.DeleteDirectory(objectId)
}

func (s *DirectoryService) GetIdMainDirectory(userID int) (int, error) {
	return s.repo.GetIdMainDirectory(userID)
}
