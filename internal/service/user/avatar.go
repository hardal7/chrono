package user

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
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
	logger.Info("Uploading avatar")
	limited := io.LimitReader(avatarFile, maxBytes)
	fileBytes, err := io.ReadAll(limited)
	if err != nil {
		logger.Error("Failed to read file", err)
		return err
	}
	if len(fileBytes) > maxBytes {
		logger.Error("File size too large")
		return err
	}
	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		logger.Error("Invalid filetype")
		return err
	}

	userID := ctx.Value(middleware.UserID).(uuid.UUID)
	err = conn.Queries.CreateAvatar(ctx, userID)
	if err != nil {
		logger.Error("Failed to query database", err)
		return err
	}
	avatar, err := conn.Queries.GetAvatarFromUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to query database", err)
		return err
	}
	err = createFile(fileBytes, avatar.ID.String())
	if err != nil {
		logger.Error("Failed to create file", err)
		return err
	}
	logger.Info("Uploaded avatar")
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
