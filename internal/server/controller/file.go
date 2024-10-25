package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"io"
	"net/http"
	"noah/internal/server/gateway"
	"noah/internal/server/request"
	"noah/internal/server/response"
	"noah/pkg/conn"
	"os"
	"strconv"
)

type FileController struct {
	gateway *gateway.Gateway
}

func NewFileController(i do.Injector) *FileController {
	return &FileController{
		gateway: do.MustInvoke[*gateway.Gateway](i),
	}
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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)

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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)
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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)
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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)
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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)
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
	_, err = f.gateway.SendCommand(uint(id), conn.Download, request.DownloadRequest{Filename: localFilename, Path: path + "/" + file.Filename}, false)
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

	result, err := f.gateway.SendCommand(uint(id), conn.FileExplorer, query, true)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}
