package user

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	maxBytes  = 1024 * 1024 * 5 // 5 MB
	dirPerm   = 0755
	filePerm  = 0644 // Don't execute the file
	directory = "uploads"
)

func UploadAvatar(ctx context.Context, avatar io.Reader) error {
	logger.Info("Uploading avatar")

	fileBytes, err := io.ReadAll(avatar)
	if len(fileBytes) > maxBytes {
		return err
	}
	if err != nil {
		return err
	}
	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return err
	}

	err = createFile(fileBytes, strconv.Itoa(ctx.Value(middleware.UserID).(int)))
	if err != nil {
		return err
	}
	logger.Info("Uploaded avatar")
	return nil
}

func createFile(fileBytes []byte, filename string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.Mkdir(directory, dirPerm)
		if err != nil {
			return err
		}
	}
	err := os.WriteFile(filepath.Join(directory, filename), fileBytes, filePerm)
	if err != nil {
		return err
	}
	return nil
}
