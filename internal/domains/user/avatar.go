package user

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hardal7/chrono/internal/middleware"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	maxBytes  = 1024 * 1024 * 5 // 5 MB
	dirPerm   = 0755
	filePerm  = 0644 // Don't execute the file
	directory = "uploads"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	logger.Info("Uploading avatar")

	fileBytes, err := io.ReadAll(r.Body)
	if len(fileBytes) > maxBytes {
		ErrLargeFile.Handle(w, err)
		return
	}
	if err != nil {
		ErrInvalidFile.Handle(w, err)
		return
	}
	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/png" {
		ErrInvalidFile.Handle(w, err)
		return
	}

	err = createFile(fileBytes, strconv.Itoa(r.Context().Value(middleware.UserID).(int)))
	if err != nil {
		e.ErrCreate.Handle(w, err, "avatar file")
	} else {
		logger.Info("Uploaded avatar")
	}
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

var ErrLargeFile = e.Error{
	InternalInfo: "File exceeds upload limit",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "File exceeds upload limit: 5MB",
}
var ErrInvalidFile = e.Error{
	InternalInfo: "Invalid filetype uploaded",
	Code:         http.StatusUnauthorized,
	ExternalInfo: "Invalid filetype uploaded, valid filetypes are: JPEG/PNG",
}
