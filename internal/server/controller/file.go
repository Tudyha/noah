package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"noah/internal/server/dto"
	"noah/internal/server/service"
	"noah/internal/server/utils"
	"noah/internal/server/vo"
	"os"
	"strconv"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

func (d *FileController) GetFileList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:   "list",
		Path: path,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", utils.ByteToString(jsonStr))

	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	var fileList []dto.FileExplorer
	err = json.Unmarshal([]byte(result), &fileList)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, fileList)
}

func (d *FileController) GetFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:   "cat",
		Path: path,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", utils.ByteToString(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (d *FileController) RenameFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body vo.DeviceFileRenamePostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:       "rename",
		Path:     body.Path,
		Filename: body.Filename,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", utils.ByteToString(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (d *FileController) DeleteFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body vo.DeviceFileDeletePostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:   "remove",
		Path: body.Path,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", string(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (d *FileController) UpdateFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body vo.DeviceFileContentPostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:          "edit",
		Path:        body.Path,
		FileContent: body.Content,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", string(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}

func (d *FileController) UploadFile(c *gin.Context) {
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

	m := make(map[string]string)
	m["path"] = path + "/" + file.Filename
	m["filename"] = localFilename
	//map转json字符串
	jsonStr, err := json.Marshal(m)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	//发送命令，让客户端来上传文件
	_, err = service.GetClientService().SendCommand(uint(id), "download", utils.ByteToString(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, "success")
}

func (d *FileController) NewDir(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body vo.DeviceFileNewDirPostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := &dto.FileExplorerQueryDto{
		Op:   "mkdir",
		Path: body.Path,
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "explorer", utils.ByteToString(jsonStr))
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, result)
}
