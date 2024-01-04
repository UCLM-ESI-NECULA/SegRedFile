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
	CreateFile(username, docID string, content []byte) (int, error)
	UpdateFile(username, docID string, content []byte) (int, error)
	DeleteFile(username, docID string) error
	GetAllUserDocs(username string) (*map[string]string, error)
}

func (fs *FileServiceImpl) GetFile(username, docID string) (string, error) {
	content, err := fs.repo.GetFile(username, docID)
	return content, err
}

func (fs *FileServiceImpl) CreateFile(username, docID string, content []byte) (int, error) {
	size, err := fs.repo.CreateFile(username, docID, content)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (fs *FileServiceImpl) UpdateFile(username, docID string, content []byte) (int, error) {
	size, err := fs.repo.UpdateFile(username, docID, content)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (fs *FileServiceImpl) DeleteFile(username, docID string) error {
	return fs.repo.DeleteFile(username, docID)
}

func (fs *FileServiceImpl) GetAllUserDocs(username string) (*map[string]string, error) {
	docs, err := fs.repo.GetAllUserDocs(username)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
