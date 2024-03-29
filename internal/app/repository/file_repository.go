package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"seg-red-file/internal/app/common"
)

type FileRepositoryImpl struct {
	// Base directory where user files are stored
	baseDir string
}

func NewFileRepository(baseDir string) *FileRepositoryImpl {
	return &FileRepositoryImpl{baseDir}
}

type FileRepository interface {
	GetFile(username, docID string) (string, error)
	CreateFile(username, docID string, content []byte) (int, error)
	UpdateFile(username, docID string, content []byte) (int, error)
	DeleteFile(username, docID string) error
	GetAllUserDocs(username string) (*map[string]string, error)
}

// UTILS

func (fr *FileRepositoryImpl) filePath(username, docID string) string {
	return filepath.Join(fr.baseDir, username, docID+".json")
}

func (fr *FileRepositoryImpl) directoryPath(username string) string {
	return filepath.Join(fr.baseDir, username)
}

func (fr *FileRepositoryImpl) ensureUserFolderExists(username string) error {
	userFolderPath := filepath.Join(fr.baseDir, username)
	return os.MkdirAll(userFolderPath, 0755)
}

//READ ONLY

func (fr *FileRepositoryImpl) GetFile(username, docID string) (string, error) {
	filePath := fr.filePath(username, docID)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", common.BadRequestError("error reading file: " + err.Error())
	}
	return string(content), nil
}

func (fr *FileRepositoryImpl) GetAllUserDocs(username string) (*map[string]string, error) {
	userFolderPath := filepath.Join(fr.baseDir, username)
	files, err := os.ReadDir(userFolderPath)
	if err != nil {
		return nil, common.BadRequestError("error reading directory: " + err.Error())
	}

	userDocs := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			docID := file.Name()
			content, err := os.ReadFile(filepath.Join(userFolderPath, docID))
			if err != nil {
				return nil, common.BadRequestError("error reading file: " + err.Error())
			}
			userDocs[docID] = string(content)
		}
	}

	return &userDocs, nil
}

//READ/WRITE

func (fr *FileRepositoryImpl) CreateFile(username, docID string, content []byte) (int, error) {
	if err := fr.ensureUserFolderExists(username); err != nil {
		return 0, fmt.Errorf("error creating user folder: %v", err)
	}

	filePath := fr.filePath(username, docID)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return 0, common.BadRequestError("file with name already exists: " + docID)
	}
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return 0, fmt.Errorf("error writing file: %v", err)
	}
	return len(content), nil
}

func (fr *FileRepositoryImpl) UpdateFile(username, docID string, content []byte) (int, error) {
	if err := fr.ensureUserFolderExists(username); err != nil {
		return 0, fmt.Errorf("error creating user folder: %v", err)
	}
	filePath := fr.filePath(username, docID)
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return 0, fmt.Errorf("error updating file: %v", err)
	}
	return len(content), nil
}

func (fr *FileRepositoryImpl) DeleteFile(username, docID string) error {
	filePath := fr.filePath(username, docID)

	// Check if file already exists
	_, existingErr := os.Stat(filePath)
	if os.IsNotExist(existingErr) {
		return common.BadRequestError("file with name: " + docID + " does not exist")
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	dirPath := fr.directoryPath(username)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	if len(files) == 0 {
		err = os.Remove(dirPath)
		if err != nil {
			return fmt.Errorf("error deleting directory: %v", err)
		}
	}
	return nil
}
