package service

import (
	"github.com/danil-vas/cloud-storage/pkg/repository"
	"time"
)

type FileService struct {
	repo repository.File
}

func NewFileService(repo repository.File) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) PathUploadFile(userId int, objectId int) (string, error) {

	return s.repo.PathUploadFile(userId, objectId)
}

func (s *FileService) AddUploadFileToUser(userId int, objectId int, originalName string, serverName string, size int64, create_time time.Time) (int, error) {
	return s.repo.AddUploadFileToUser(userId, objectId, originalName, serverName, size, create_time)
}

func (s *FileService) OriginalFileName(objectId int) (string, error) {
	return s.repo.OriginalFileName(objectId)
}

func (s *FileService) DeleteFile(userId int, objectId int) error {
	return s.repo.DeleteFile(userId, objectId)
}

func (s *FileService) GetAvailableMemory(userId int) (int, error) {
	return s.repo.GetAvailableMemory(userId)
}

func (s *FileService) CheckAccessToObject(userId int, objectId int) (bool, error) {
	return s.repo.CheckAccessToObject(userId, objectId)
}

func (s *FileService) GetTypeObject(objectId int) (string, error) {
	return s.repo.GetTypeObject(objectId)
}
func (s *FileService) OriginalFileNameThroughServerName(serverName string) (string, error) {
	return s.repo.OriginalFileNameThroughServerName(serverName)
}
