package controller

import (
	"io"
	"net/http"
	"seg-red-file/internal/app/common"
	"seg-red-file/internal/app/repository"
	"seg-red-file/internal/app/service"

	"github.com/gin-gonic/gin"
)

type FileControllerImpl struct {
	svc service.FileService
}

func NewFileController() *FileControllerImpl {
	return &FileControllerImpl{
		svc: service.NewFileService(repository.NewFileRepository(""))}
}

type FileController interface {
	GetFile(c *gin.Context)
	CreateFile(c *gin.Context)
	UpdateFile(c *gin.Context)
	DeleteFile(c *gin.Context)
	GetAllUserDocs(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (fc *FileControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:username/:doc_id", fc.GetFile)
	router.POST("/:username/:doc_id", fc.CreateFile)
	router.PUT("/:username/:doc_id", fc.UpdateFile)
	router.DELETE("/:username/:doc_id", fc.DeleteFile)
	router.GET("/:username/_all_docs", fc.GetAllUserDocs)
}

func (fc *FileControllerImpl) GetFile(c *gin.Context) {
	username, docID := checkParams(c)
	content, err := fc.svc.GetFile(username, docID)
	if err != nil {
		common.NewAPIError(c, http.StatusNotFound, err, "file not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (fc *FileControllerImpl) CreateFile(c *gin.Context) {
	username, docID := checkParams(c)

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	size := fc.svc.CreateFile(username, docID, requestBody)
	c.JSON(http.StatusOK, gin.H{"size": size})
}

func (fc *FileControllerImpl) UpdateFile(c *gin.Context) {
	username, docID := checkParams(c)
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	size := fc.svc.UpdateFile(username, docID, requestBody)
	c.JSON(http.StatusOK, gin.H{"size": size})
}

func (fc *FileControllerImpl) DeleteFile(c *gin.Context) {
	username, docID := checkParams(c)
	err := fc.svc.DeleteFile(username, docID)
	if err != nil {
		common.NewAPIError(c, http.StatusNotFound, err, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (fc *FileControllerImpl) GetAllUserDocs(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		common.NewAPIError(c, http.StatusBadRequest, nil, "username cannot be empty")
		return
	}
	docs := fc.svc.GetAllUserDocs(username)
	if docs == nil {
		docs = make(map[string]string)
	}
	c.JSON(http.StatusOK, docs)
}

// checkParams checks if the username and docID are valid
func checkParams(c *gin.Context) (string, string) {
	username := c.Param("username")
	docID := c.Param("doc_id")
	if username == "" || docID == "" {
		common.NewAPIError(c, http.StatusBadRequest, nil, "invalid input parameters")
	}
	return username, docID
}
