package explorer

import (
	"errors"
	"fmt"
	"log"
	"noah/client/app/utils"
	"os"
	"path/filepath"
	"strings"

	"noah/client/app/entitie"
	"noah/client/app/service"
)

type Service struct{}

func NewService() service.FileExplorer {
	return &Service{}
}

// GetFileExplorer 返回指定路径下的所有文件和目录信息。
func (s Service) GetFileExplorer(path string) ([]entitie.FileExplorer, error) {
	if !isValidPath(path) {
		return nil, errors.New("invalid path")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []entitie.FileExplorer
	for _, fileEntry := range files {
		filePath := filepath.Join(path, fileEntry.Name())
		fileInfo, err := fileEntry.Info()
		if err != nil {
			// 忽略无法获取文件信息的情况
			continue
		}
		fileType := getFileType(fileEntry)
		result = append(result, entitie.FileExplorer{
			Filename: fileEntry.Name(),
			ModTime:  fileInfo.ModTime(),
			Path:     filePath,
			Type:     uint8(fileType),
		})
	}
	return result, nil
}

// getFileType 根据文件类型返回对应的类型值。
func getFileType(f os.DirEntry) uint8 {
	if f.IsDir() {
		return 2
	}
	fileInfo, _ := f.Info()
	switch fileInfo.Mode() & os.ModeType {
	case os.ModeSymlink:
		return 3
	case os.ModeDevice:
		return 4
	default:
		return 1
	}
}

// isValidPath 验证路径是否有效。
func isValidPath(path string) bool {
	return !strings.Contains(path, "..")
}

func (s Service) ReadFile(path string) ([]byte, error) {
	if !isValidPath(path) {
		return nil, errors.New("invalid path")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Rename 重命名文件或文件夹
func (s Service) Rename(path string, newFilename string) error {
	// 参数验证
	if path == "" || newFilename == "" {
		return fmt.Errorf("path or newFilename is empty")
	}

	// 安全性检查
	if filepath.IsAbs(newFilename) {
		return fmt.Errorf("newFilename should not be an absolute path")
	}

	// 获取旧文件名
	oldFilename := filepath.Base(path)
	if oldFilename == "" {
		return fmt.Errorf("failed to extract oldFilename from path")
	}

	// 拼接完整路径
	dirPath := filepath.Dir(path)
	oldPath := filepath.Join(dirPath, oldFilename)
	newPath := filepath.Join(dirPath, newFilename)

	// 尝试重命名
	err := os.Rename(oldPath, newPath)
	if err != nil {
		// 更详细的错误处理
		if os.IsNotExist(err) {
			return fmt.Errorf("source path does not exist: %v", err)
		} else if os.IsExist(err) {
			return fmt.Errorf("destination path already exists: %v", err)
		} else {
			return fmt.Errorf("failed to rename file or directory: %v", err)
		}
	}

	return nil
}

// Remove 删除指定路径及其子目录下的所有文件和目录
func (s Service) Remove(path string) error {
	// 参数校验
	if path == "" {
		return fmt.Errorf("invalid path: path is empty")
	}

	// 执行删除操作
	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("Failed to remove path %s: %v", path, err)
		return fmt.Errorf("failed to remove path %s: %w", path, err)
	}

	log.Printf("Successfully removed path %s", path)
	return nil
}

// WriteFile 将给定的内容写入指定的文件路径。
// 如果路径不存在或者写入失败，将返回错误。
func (s Service) WriteFile(path string, content []byte) error {
	return utils.WriteFile(path, content)
}

func (s Service) MkDir(path string) error {
	return os.Mkdir(path, 0755)
}
