package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/mailgun/mailgun-go/v5"
)

const (
	otpExpiration = time.Minute * 5
	mailTimeout   = time.Second * 10
)

func RequestPasswordReset(ctx context.Context, r dto.RequestUserPasswordResetRequest) error {
	if r.Email == "" && r.Username == "" {
		return errors.New("Both email and username fields cannot be empty")
	}

	var err error
	var user query.User
	if r.Email == "" {
		user, err = db.Queries.GetUserByUsername(ctx, r.Username)
		if err != nil {
			return fmt.Errorf("Failed to get user by username: %w: %w", db.ErrRunQuery, err)
		}
	} else {
		user, err = db.Queries.GetUserByEmail(ctx, r.Email)
		if err != nil {
			return fmt.Errorf("Failed to get user by username: %w: %w", db.ErrRunQuery, err)
		}
	}

	auth.AsUserID(ctx, user.ID)

	err = sendResetEmail(ctx, user.Email)

	return err
}

func PasswordReset(ctx context.Context, otp string, r dto.UserPasswordResetRequest) error {
	hashedOTP := auth.HashToken(otp, []byte(config.App.HashSecret))
	token, err := db.Queries.GetOTPToken(ctx, hashedOTP)
	if err != nil {
		return fmt.Errorf("Failed to retrieve OTP token: %w: %w", db.ErrRunQuery, err)
	}

	if !token.Expiry.After(time.Now()) {
		return errors.New("OTP token has expired")
	}

	ctx = auth.AsUserID(ctx, token.UserID)
	err = EditAccount(ctx, dto.EditUserAccountRequest{NewPassword: r.NewPassword})
	if err != nil {
		return err
	}

	err = db.Queries.DeleteOTPToken(ctx, token.ID)
	if err != nil {
		return fmt.Errorf("Failed to invalidate consumed OTP token: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}

func sendResetEmail(ctx context.Context, email string) error {
	mg := mailgun.NewMailgun(config.App.MailAPIKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		logger.Warn("Failed to set API Base for email")
		return err
	}

	token, err := generateOTPToken(ctx)
	if err != nil {
		return err
	}
	domain := config.App.Domain

	recipient := email
	sender := config.App.MailAddress
	subject := "Password Reset"
	body := "Click this link to reset your password: https://" + domain + "/password-reset?otp=" + token

	data, err := os.ReadFile("static/mail/password-reset.html")
	if err != nil {
		return err
	}
	html := string(data)
	html = strings.ReplaceAll(html, "RESET_TOKEN", token)
	html = strings.ReplaceAll(html, "DOMAIN_NAME", domain)

	message := mailgun.NewMessage(domain, sender, subject, body, recipient)
	message.SetHTML(html)

	ctx, cancel := context.WithTimeout(context.Background(), mailTimeout)
	defer cancel()

	_, err = mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func generateOTPToken(ctx context.Context) (string, error) {
	var otp string
	token, err := auth.GenerateToken()
	if err != nil {
		return otp, fmt.Errorf("Failed to generate token: %w", err)
	}

	hashedToken := auth.HashToken(token, []byte(config.App.HashSecret))

	err = db.Queries.CreateOTPToken(ctx, query.CreateOTPTokenParams{
		UserID: auth.UserID(ctx),
		Expiry: time.Now().Add(otpExpiration),
		Hash:   hashedToken,
	})
	if err != nil {
		return otp, fmt.Errorf("Failed to create session token: %w: %w", db.ErrRunQuery, err)
	}

	return token, nil
}
