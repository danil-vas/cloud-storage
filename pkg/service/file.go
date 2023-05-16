package service

import "github.com/danil-vas/cloud-storage/pkg/repository"

type FileService struct {
	repo repository.File
}

func NewFileService(repo repository.File) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) PathUploadFile(userId int, objectId int) (string, error) {

	return s.repo.PathUploadFile(userId, objectId)
}
