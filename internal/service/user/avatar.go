package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	maxBytes = 1024 * 1024 * 5 // 5 MB
	filePerm = 0o644           // Don't execute the file

	AvatarDirectory        = "/srv/avatars"
	DefaultAvatarDirectory = "default"
	defaultAvatarsNum      = 15
)

// TODO: Sanitize Image
func UploadAvatar(ctx context.Context, avatarFile io.Reader) error {
	userID := auth.UserID(ctx)

	err := DeleteAvatar(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete previous avatar: %w", err)
	}

	limited := io.LimitReader(avatarFile, maxBytes)
	fileBytes, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("Failed to read file: %w", err)
	}

	if len(fileBytes) > maxBytes {
		return errors.New("File size too large")
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return errors.New("Invalid filetype")
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

func DeleteAvatar(ctx context.Context) error {
	userID := auth.UserID(ctx)

	path := filepath.Join(AvatarDirectory, userID.String())
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Failed to delete user avatar: %w", err)
	}

	err = InitAvatar(ctx)
	if err != nil {
		logger.Warn("Failed to initialize user avatar")
	}

	return nil
}

func InitAvatar(ctx context.Context) error {
	userID := auth.UserID(ctx)
	n, err := rand.Int(rand.Reader, big.NewInt(defaultAvatarsNum))
	if err != nil {
		logger.Warn("Failed to generate random number")
		return err
	}
	randomAvatar := n.String()
	avatarPath := filepath.Join(DefaultAvatarDirectory, randomAvatar)

	err = createSymlink(avatarPath, userID.String())

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
