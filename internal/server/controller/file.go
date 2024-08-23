package controller

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"noah/internal/server/dto"
	"noah/internal/server/service"
	"noah/internal/server/utils"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

	cmd := "ls -l " + path

	result, err := service.GetClientService().SendCommand(uint(id), cmd, "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if result == "No content." {
		Fail(c, http.StatusBadRequest, "No content.")
		return
	}
	var data []dto.DeviceFileDto
	for i, line := range strings.Split(result, "\n") {
		if i == 0 {
			continue
		}
		if line == "" {
			continue
		}
		name, fType, err := parseLsL(line)
		if err != nil {
			Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		var pPath string
		if path == "/" {
			pPath = path
		} else {
			pPath = path + "/"
		}
		data = append(data, dto.DeviceFileDto{
			Type: fType,
			Name: name,
			Path: pPath + name,
		})
	}

	Success(c, data)
}

func (d *FileController) GetFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	path := c.Query("path")
	if path == "" {
		Fail(c, http.StatusBadRequest, "path is required")
		return
	}

	result, err := service.GetClientService().SendCommand(uint(id), "cat "+path, "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if result == "No content." {
		Success(c, "")
		return
	}

	Success(c, result)
}

type FileRenamePostVo struct {
	Name string `json:"name" binding:"required"`
	Path string `json:"path" binding:"required"`
}

func (d *FileController) RenameFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body FileRenamePostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	//去掉/分割的最后一个元素
	pPath := strings.Join(strings.Split(body.Path, "/")[:len(strings.Split(body.Path, "/"))-1], "/")

	result, err := service.GetClientService().SendCommand(uint(id), "mv "+body.Path+" "+pPath+"/"+body.Name, "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if result == "No content." {
		Success(c, "")
		return
	}

	Success(c, result)
}

type FileDeletePostVo struct {
	Path string `json:"path" binding:"required"`
	Type int8   `json:"type" binding:"required"`
}

func (d *FileController) DeleteFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body FileDeletePostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	cmd := "rm -rf " + body.Path
	if body.Type == 1 {
		cmd = "rm -f " + body.Path
	}

	result, err := service.GetClientService().SendCommand(uint(id), cmd, "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if result == "No content." {
		Success(c, "")
		return
	}

	Success(c, result)
}

func parseLsL(line string) (string, int8, error) {
	name := strings.Fields(line)[8]

	re := regexp.MustCompile(`^([d-])`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		if len(strings.Fields(line)) < 11 {
			return name, 3, nil
		} else if len(strings.Fields(line)) >= 11 {
			return name + strings.Fields(line)[9] + strings.Fields(line)[10], 3, nil
		}
	}

	isDir := strings.HasPrefix(matches[1], "d")
	if isDir {
		return name, 2, nil
	}
	return name, 1, nil
}

type FileContentPostVo struct {
	Content string `json:"content" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

func (d *FileController) UpdateFileContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body FileContentPostVo
	err := c.BindJSON(&body)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = service.GetClientService().SendCommand(uint(id), fmt.Sprintf("echo %s > %s", body.Content, body.Path), "")
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, "success")
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
