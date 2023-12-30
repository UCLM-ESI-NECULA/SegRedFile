package service

import (
	"seg-red-file/internal/app/repository"
)

type FileServiceImpl struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) *FileServiceImpl {
	return &FileServiceImpl{repo}
}

type FileService interface {
	GetFile(username, docID string) (string, error)
	CreateFile(username, docID string, content []byte) int
	UpdateFile(username, docID string, content []byte) int
	DeleteFile(username, docID string) error
	GetAllUserDocs(username string) map[string]string
}

func (fs *FileServiceImpl) GetFile(username, docID string) (string, error) {
	content, err := fs.repo.GetFile(username, docID)
	return content, err
}

func (fs *FileServiceImpl) CreateFile(username, docID string, content []byte) int {
	size, err := fs.repo.CreateFile(username, docID, content)
	if err != nil {
		// Handle specific errors (e.g., write errors, permission issues)
		return 0
	}
	return size
}

func (fs *FileServiceImpl) UpdateFile(username, docID string, content []byte) int {
	size, _ := fs.repo.UpdateFile(username, docID, content)
	return size
}

func (fs *FileServiceImpl) DeleteFile(username, docID string) error {
	return fs.repo.DeleteFile(username, docID)
}

func (fs *FileServiceImpl) GetAllUserDocs(username string) map[string]string {
	docs, _ := fs.repo.GetAllUserDocs(username)
	return docs
}
