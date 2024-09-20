package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"noah/internal/server/enum"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/internal/server/service"
	"os"
	"strconv"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

func (f FileController) GetFileList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:   "list",
		Path: path,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)

	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	var fileList []response.GetFileExplorerRes
	err = json.Unmarshal([]byte(result), &fileList)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, fileList)
}

func (f FileController) GetFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:   "cat",
		Path: path,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (f FileController) RenameFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileRenameReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:       "rename",
		Path:     body.Path,
		Filename: body.Filename,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (f FileController) DeleteFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileDeleteReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:   "remove",
		Path: body.Path,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (f FileController) UpdateFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileContentReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:          "edit",
		Path:        body.Path,
		FileContent: body.Content,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (f FileController) UploadFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	path := c.PostForm("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		Fail(c, http.StatusBadRequest, fmt.Sprintf("Error uploading file: %s", err.Error()))
		return
	}

	localFilename := uuid.New().String()
	saveFilePath := "temp/" + localFilename

	out, err := os.Create(saveFilePath)
	if err != nil {
		Fail(c, http.StatusInternalServerError, fmt.Sprintf("Error creating file: %s", err.Error()))
		return
	}
	defer out.Close()

	src, err := file.Open()
	if err != nil {
		Fail(c, http.StatusInternalServerError, fmt.Sprintf("Error opening file: %s", err.Error()))
		return
	}
	defer src.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		Fail(c, http.StatusInternalServerError, fmt.Sprintf("Error saving file: %s", err.Error()))
		return
	}

	//发送命令，让客户端来上传文件
	_, err = service.GetChannelService().SendCommand(uint(id), enum.MessageTypeDownload, request.DownloadRequest{Filename: localFilename, Path: path + "/" + file.Filename})
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, "success")
}

func (f FileController) NewDir(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileNewDirReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &request.GetFileExplorerQueryReq{
		Op:   "mkdir",
		Path: body.Path,
	}

	result, err := service.GetChannelService().SendCommand(uint(id), enum.MessageTypeFileExplorer, query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}
