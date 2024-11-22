package controller

import (
	"encoding/json"
	"io"
	"noah/internal/server/gateway"
	"noah/pkg/errcode"
	"noah/pkg/mux/message"
	"os"
	"strconv"

	"noah/pkg/request"

	"noah/pkg/enum"

	"noah/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type FileController struct {
	gateway *gateway.Gateway
}

func NewFileController(i do.Injector) (FileController, error) {
	return FileController{
		gateway: do.MustInvoke[*gateway.Gateway](i),
	}, nil
}

func (f FileController) GetFileList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:   "list",
		Path: path,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)

	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}
	var fileList []response.GetFileExplorerRes
	err = json.Unmarshal([]byte(result), &fileList)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, fileList)
}

func (f FileController) GetFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:   "cat",
		Path: path,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, result)
}

func (f FileController) RenameFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileRenameReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:       "rename",
		Path:     body.Path,
		Filename: body.Filename,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, result)
}

func (f FileController) DeleteFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileDeleteReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:   "remove",
		Path: body.Path,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, result)
}

func (f FileController) UpdateFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileContentReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:          "edit",
		Path:        body.Path,
		FileContent: body.Content,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, result)
}

func (f FileController) UploadFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	path := c.PostForm("path")
	if path == "" {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	localFilename := uuid.New().String()
	saveFilePath := "temp/" + localFilename

	out, err := os.Create(saveFilePath)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}
	defer out.Close()

	src, err := file.Open()
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}
	defer src.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	//发送命令，让客户端来下载文件
	_, err = f.gateway.SendCommand(uint(id), enum.Download, message.DownloadReq{Filename: localFilename, Path: path + "/" + file.Filename}, false)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, "success")
}

func (f FileController) NewDir(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body request.ClientFileNewDirReq
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, errcode.ErrInvalidParameter)
		return
	}

	query := message.FileExplorerReq{
		Op:   "mkdir",
		Path: body.Path,
	}

	result, err := f.gateway.SendCommand(uint(id), enum.FileExplorer, query, true)
	if err != nil {
		Fail(c, errcode.ErrInternalError)
		return
	}

	Success(c, result)
}
