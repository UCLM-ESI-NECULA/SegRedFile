package controller

import (
	"io"
	"net/http"
	"seg-red-file/internal/app/common"
	"seg-red-file/internal/app/repository"
	"seg-red-file/internal/app/service"
	"seg-red-file/internal/dao"

	"github.com/gin-gonic/gin"
)

type FileControllerImpl struct {
	svc service.FileService
}

func NewFileController(g *gin.RouterGroup) *FileControllerImpl {
	c := &FileControllerImpl{
		svc: service.NewFileService(repository.NewFileRepository(""))}
	c.RegisterRoutes(g)
	return c
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
	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	content, err := fc.svc.GetFile(username, docID)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.FileContent{Content: content})
}

func (fc *FileControllerImpl) CreateFile(c *gin.Context) {
	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	//ReadBody
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ForwardError(c, common.BadRequestError("invalid request body"))
		return
	}

	size, err := fc.svc.CreateFile(username, docID, requestBody)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.FileSize{Size: size})
}

func (fc *FileControllerImpl) UpdateFile(c *gin.Context) {
	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	//ReadBody
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ForwardError(c, common.BadRequestError("invalid request body"))
		return
	}

	size, err := fc.svc.UpdateFile(username, docID, requestBody)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.FileSize{Size: size})
}

func (fc *FileControllerImpl) DeleteFile(c *gin.Context) {
	// Check username and docID
	username, docID, apiErr := checkParams(c)
	if apiErr != nil {
		common.ForwardError(c, apiErr)
		return
	}

	err := fc.svc.DeleteFile(username, docID)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (fc *FileControllerImpl) GetAllUserDocs(c *gin.Context) {
	// Check username
	username := c.Param("username")
	if username == "" {
		common.ForwardError(c, common.EmptyParamsError("username"))
		return
	}

	docs, errors := fc.svc.GetAllUserDocs(username)
	if errors != nil {
		common.HandleError(c, errors)
		return
	}
	if docs == nil {
		m := make(map[string]string)
		docs = &m
	}
	c.JSON(http.StatusOK, docs)
}

// checkParams checks if the username and docID are valid
func checkParams(c *gin.Context) (string, string, *common.APIError) {
	username := c.Param("username")
	if username == "" {
		return "", "", common.EmptyParamsError("username")
	}
	docID := c.Param("doc_id")
	if docID == "" {
		return "", "", common.EmptyParamsError("doc_id")
	}
	return username, docID, nil
}
