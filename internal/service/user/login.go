package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, r dto.LoginUserRequest) (http.Cookie, error) {
	var err error
	var u query.User
	var cookie http.Cookie

	if r.Username != "" {
		u, err = db.Queries.GetUserByUsername(ctx, r.Username)
	} else if r.Email != "" {
		u, err = db.Queries.GetUserByEmail(ctx, r.Email)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return cookie, fmt.Errorf("user not found: %w", db.ErrNotFound)
	} else if err != nil {
		return cookie, fmt.Errorf("get user: %w: %w", db.ErrRunQuery, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return cookie, fmt.Errorf("wrong password: %w", err)
	} else if err != nil {
		return cookie, fmt.Errorf("hash password: %w", err)
	}

	ctx = requestctx.AsUserID(ctx, u.ID)
	cookie, err = auth.GenerateCookie(ctx)
	return cookie, err
}
