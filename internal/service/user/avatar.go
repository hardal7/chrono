package user

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hardal7/chrono/internal/middleware"
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
	userID := middleware.UserID(ctx)

	limited := io.LimitReader(avatarFile, maxBytes)
	fileBytes, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("Failed to read file: %w", err)
	}

	if len(fileBytes) > maxBytes {
		return fmt.Errorf("File size too large")
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return fmt.Errorf("Invalid filetype")
	}

	err = createFile(fileBytes, userID.String())
	if err != nil {
		return fmt.Errorf("Failed to create file: %w", err)
	}

	return nil
}

func createFile(fileBytes []byte, filename string) error {
	path := filepath.Join(AvatarDirectory, filename)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Failed to delete file: %w", err)
	}

	err = os.WriteFile(path, fileBytes, filePerm)
	if err != nil {
		return fmt.Errorf("Failed to write file: %w", err)
	}

	return nil
}

func InitAvatar(ctx context.Context) error {
	userID := middleware.UserID(ctx)

	randomAvatar := strconv.Itoa(rand.IntN(defaultAvatarsNum))
	avatarPath := filepath.Join(DefaultAvatarDirectory, randomAvatar)

	err := createSymlink(avatarPath, userID.String())

	return err
}

func createSymlink(source, filename string) error {
	path := filepath.Join(AvatarDirectory, filename)
	err := os.Symlink(source, path)
	if err != nil {
		return fmt.Errorf("Failed to create symlink: %w", err)
	}

	return nil
}
