package user

import (
	"context"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	maxBytes = 1024 * 1024 * 5 // 5 MB
	filePerm = 0644            // Don't execute the file

	AvatarDirectory        = "/srv/avatars"
	DefaultAvatarDirectory = "default"
	defaultAvatarsNum      = 15
)

// TODO: Sanitize Image
func UploadAvatar(ctx context.Context, avatarFile io.Reader) error {
	logger.Debug("Uploading user avatar")
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

	logger.Debug("Uploaded user avatar")
	return nil
}

func createFile(fileBytes []byte, filename string) error {
	path := filepath.Join(AvatarDirectory, filename)
	err := os.Remove(path)
	if err != nil {
		logger.Debug("Failed to delete file", err)
		return err
	}

	err = os.WriteFile(path, fileBytes, filePerm)
	if err != nil {
		logger.Debug("Failed to write file", err)
		return err
	}

	return nil
}

func InitAvatar(ctx context.Context) error {
	logger.Debug("Initializing user avatar")

	randomAvatar := strconv.Itoa(rand.IntN(defaultAvatarsNum))
	avatarPath := filepath.Join(DefaultAvatarDirectory, randomAvatar)
	userID := ctx.Value(middleware.UserID).(uuid.UUID).String()

	err := createSymlink(avatarPath, userID)

	logger.Debug("Initialized user avatar", "id", randomAvatar)
	return err
}

func createSymlink(source, filename string) error {
	path := filepath.Join(AvatarDirectory, filename)
	err := os.Symlink(source, path)
	if err != nil {
		logger.Debug("Failed to create symlink", err)
		return err
	}
	return nil
}
