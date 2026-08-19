package user

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	maxBytes        = 1024 * 1024 * 5 // 5 MB
	dirPerm         = 0755
	filePerm        = 0644 // Don't execute the file
	AvatarDirectory = "avatars"
)

func UploadAvatar(ctx context.Context, avatarFile io.Reader) error {
	logger.Debug("Uploading avatar")
	limited := io.LimitReader(avatarFile, maxBytes)
	fileBytes, err := io.ReadAll(limited)
	if err != nil {
		logger.Debug("Failed to read file", err)
		return err
	}
	if len(fileBytes) > maxBytes {
		logger.Debug("File size too large")
		return err
	}
	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		logger.Debug("Invalid filetype")
		return err
	}

	userID := ctx.Value(middleware.UserID).(uuid.UUID).String()
	err = createFile(fileBytes, userID)
	if err != nil {
		logger.Debug("Failed to create file", err)
		return err
	}
	logger.Debug("Uploaded avatar")
	return nil
}

func createFile(fileBytes []byte, filename string) error {
	if _, err := os.Stat(AvatarDirectory); os.IsNotExist(err) {
		err = os.Mkdir(AvatarDirectory, dirPerm)
		if err != nil {
			return err
		}
	}
	err := os.WriteFile(filepath.Join(AvatarDirectory, filename), fileBytes, filePerm)
	if err != nil {
		return err
	}
	return nil
}
