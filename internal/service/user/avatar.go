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
	"uuid"

	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
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
	userID := requestctx.GetUserID(ctx)

	err := DeleteAvatar(ctx)
	if err != nil {
		return fmt.Errorf("delete previous avatar: %w", err)
	}

	limited := io.LimitReader(avatarFile, maxBytes)
	fileBytes, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if len(fileBytes) > maxBytes {
		return errors.New("file size too large")
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		return fmt.Errorf("invalid filetype: %q", filetype)
	}

	err = createFile(fileBytes, userID.String())
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	return nil
}

func createFile(fileBytes []byte, filename string) error {
	path := filepath.Join(AvatarDirectory, filename)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	err = os.WriteFile(path, fileBytes, filePerm)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func DeleteAvatar(ctx context.Context) error {
	userID := requestctx.GetUserID(ctx)

	path := filepath.Join(AvatarDirectory, userID.String())
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("delete user avatar: %w", err)
	}

	err = InitAvatar(ctx)
	if err != nil {
		logger.Err(err).Warn("initialize user avatar")
	}

	return nil
}

func InitAvatar(ctx context.Context) error {
	userID := requestctx.GetUserID(ctx)
	n, err := rand.Int(rand.Reader, big.NewInt(defaultAvatarsNum))
	if err != nil {
		logger.Err(err).Warn("generate random number")
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
		return fmt.Errorf("create symlink: %w", err)
	}

	return nil
}

func GetAvatarPath(userID uuid.UUID) string {
	return filepath.Join(config.AvatarEndpoint, userID.String())
}
